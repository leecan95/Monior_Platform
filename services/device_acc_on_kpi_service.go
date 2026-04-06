package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	deviceKpiAccOnName      = "ACC_ON_OVER_10H"
	deviceKpiAccOffName     = "ACC_OFF_OVER_72H"
	deviceKpiVehicleMknName = "VEHICLE_MKN_OVER_72H"
	deviceKpiVehicleMknNow  = "VEHICLE_MKN_CURRENT"
	deviceKpiVehicleBadGps  = "VEHICLE_BADGPS_OVER_ONLINE_IN_DAY"
	deviceKpiBadWriteMemory = "BAD_WRITE_MEMORY"
	deviceKpiAccOnCycle21   = "ACC_ON_OVER_10H_CYCLE_21"
	deviceKpiAccOffCycle21  = "ACC_OFF_OVER_72H_CYCLE_21"
	deviceKpiMknCycle21     = "VEHICLE_MKN_OVER_72H_CYCLE_21"
	deviceDailyStatusCursor = "device_daily_status_plate_snapshot"
	deviceDailyStatusReset  = "device_daily_status_day_reset"
	deviceKpiInsertCursor   = "device_kpi_hourly_insert"
	deviceKpiDetailDayReset = "device_kpi_vehicle_detail_day_reset"

	statusOffline = "offline"
	statusBadGPS  = "badgps"
)

var onlineStatusValues = []string{"badgps", "overspeed", "park", "run", "stop"}

// StartDeviceAccOnKpiJob runs sampling every 5 minutes; inserting device_kpi is throttled to hourly.
func StartDeviceAccOnKpiJob() {
	go func() {
		runDeviceAccOnKpiSnapshot()

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			runDeviceAccOnKpiSnapshot()
		}
	}()
}

func runDeviceAccOnKpiSnapshot() {
	cfg := model.LoadDBConfig()
	if err := config.ValidateDBConfig(cfg); err != nil {
		log.Printf("[ERROR] device acc job skipped (db config): %v", err)
		return
	}

	timestamp := time.Now().UTC()
	sampledHourUTC := sampledHourFromUTC(timestamp, resolveKpiLocation())

	matchedMkn, totalMkn, err := countVehicleByFieldCondition("status", "offline", 72*time.Hour)
	if err != nil {
		log.Printf("[ERROR] device acc job: count VEHICLE_MKN_OVER_72H from mongo: %v", err)
		return
	}

	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] device acc job: connect trans db: %v", err)
		return
	}
	if err := db.EnsureDeviceKPITable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi table: %v", err)
		return
	}
	if err := db.EnsureDeviceKPIVehicleDetailTable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi_vehicle_detail table: %v", err)
		return
	}
	if err := db.EnsureDeviceKPIVehicleDetailSnapshotTable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi_vehicle_detail_snapshot table: %v", err)
		return
	}
	// Run day-boundary snapshot/reset independently from hourly KPI insert gate.
	if err := resetDeviceKPIVehicleDetailFlagsOnNewDay(db, timestamp); err != nil {
		log.Printf("[ERROR] device acc job: reset daily flags in device_kpi_vehicle_detail: %v", err)
		return
	}

	accStateRows, err := fetchVehicleACCStateSnapshot(timestamp)
	if err != nil {
		log.Printf("[ERROR] device acc job: fetch ACC state snapshot from mongo: %v", err)
		return
	}
	if err := db.SyncDeviceVehicleACCState(accStateRows, timestamp, sampledHourUTC); err != nil {
		log.Printf("[ERROR] device acc job: sync ACC state continuity: %v", err)
		return
	}
	matchedAccOn, matchedAccOff, err := db.CountDeviceVehicleACCViolations(sampledHourUTC)
	if err != nil {
		log.Printf("[ERROR] device acc job: count ACC_ON/ACC_OFF from postgres state: %v", err)
		return
	}
	totalAccOn := totalMkn
	totalAccOff := totalMkn

	matchedMknCurrent, matchedBadGps, matchedBadMemory, totalOnlineInDay, err := syncDailyVehicleStatusAndGetCounts(db)
	if err != nil {
		log.Printf("[ERROR] device acc job: sync daily status for VEHICLE_MKN_CURRENT and VEHICLE_BADGPS_OVER_ONLINE_IN_DAY: %v", err)
		return
	}

	cycleAccOnMatched, cycleAccOffMatched, cycleMknMatched, cycleTotalVehicles, cycleStartTs, cycleErr := countDistinctVehicleKpiInCycle21Window()
	if cycleErr != nil {
		log.Printf("[ERROR] device acc job: count cycle-21 KPI from mongo: %v", cycleErr)
	}

	canInsertDeviceKpi, insertWait, err := shouldInsertDeviceKpiSnapshot(db, timestamp)
	if err != nil {
		log.Printf("[ERROR] device acc job: decide hourly insert for device_kpi: %v", err)
		return
	}
	if !canInsertDeviceKpi {
		log.Printf("[INFO] device acc job: sampled mongo, skip device_kpi insert (next write after %s)", insertWait.UTC().Format(time.RFC3339))
		return
	}
	mknRows, err := fetchVehicleKpiDetailSnapshot(timestamp)
	if err != nil {
		log.Printf("[ERROR] device acc job: fetch hourly vehicle detail snapshot from mongo: %v", err)
		return
	}
	accOnRows, err := db.ListDeviceVehicleACCOnViolationRows(sampledHourUTC)
	if err != nil {
		log.Printf("[ERROR] device acc job: list ACC_ON vehicle detail rows from postgres state: %v", err)
		return
	}
	accOffRows, err := db.ListDeviceVehicleACCOffViolationRows(sampledHourUTC)
	if err != nil {
		log.Printf("[ERROR] device acc job: list ACC_OFF vehicle detail rows from postgres state: %v", err)
		return
	}
	detailRows := make([]model.DeviceKPIVehicleDetail, 0, len(accOnRows)+len(accOffRows)+len(mknRows))
	detailRows = append(detailRows, accOnRows...)
	detailRows = append(detailRows, accOffRows...)
	detailRows = append(detailRows, mknRows...)
	for i := range detailRows {
		detailRows[i].SampledAt = timestamp
		detailRows[i].SampledHour = sampledHourUTC
	}
	if err := db.UpsertDeviceKPIVehicleDetails(detailRows, timestamp, sampledHourUTC); err != nil {
		log.Printf("[ERROR] device acc job: upsert device_kpi_vehicle_detail rows: %v", err)
		return
	}
	log.Printf("[INFO] device acc job: upserted device_kpi_vehicle_detail rows=%d", len(detailRows))

	insertFailed := false

	accOnRow := model.DeviceKPI{
		Name:               deviceKpiAccOnName,
		MatchedDeviceCount: matchedAccOn,
		TotalDeviceCount:   totalAccOn,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(accOnRow); err != nil {
		log.Printf("[ERROR] device acc job: insert ACC_ON_OVER_10H row: %v", err)
		insertFailed = true
	}

	accOffRow := model.DeviceKPI{
		Name:               deviceKpiAccOffName,
		MatchedDeviceCount: matchedAccOff,
		TotalDeviceCount:   totalAccOff,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(accOffRow); err != nil {
		log.Printf("[ERROR] device acc job: insert ACC_OFF_OVER_72H row: %v", err)
		insertFailed = true
	}

	mknRow := model.DeviceKPI{
		Name:               deviceKpiVehicleMknName,
		MatchedDeviceCount: matchedMkn,
		TotalDeviceCount:   totalMkn,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(mknRow); err != nil {
		log.Printf("[ERROR] device acc job: insert VEHICLE_MKN_OVER_72H row: %v", err)
		insertFailed = true
	}

	mknCurrentRow := model.DeviceKPI{
		Name:               deviceKpiVehicleMknNow,
		MatchedDeviceCount: matchedMknCurrent,
		TotalDeviceCount:   totalOnlineInDay,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(mknCurrentRow); err != nil {
		log.Printf("[ERROR] device acc job: insert VEHICLE_MKN_CURRENT row: %v", err)
		insertFailed = true
	}

	badGpsRow := model.DeviceKPI{
		Name:               deviceKpiVehicleBadGps,
		MatchedDeviceCount: matchedBadGps,
		TotalDeviceCount:   totalOnlineInDay,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(badGpsRow); err != nil {
		log.Printf("[ERROR] device acc job: insert VEHICLE_BADGPS_OVER_ONLINE_IN_DAY row: %v", err)
		insertFailed = true
	}

	badWriteMemoryRow := model.DeviceKPI{
		Name:               deviceKpiBadWriteMemory,
		MatchedDeviceCount: matchedBadMemory,
		TotalDeviceCount:   totalOnlineInDay,
		Timestamp:          timestamp,
	}
	if err := db.InsertDeviceKPI(badWriteMemoryRow); err != nil {
		log.Printf("[ERROR] device acc job: insert BAD_WRITE_MEMORY row: %v", err)
		insertFailed = true
	}

	if cycleErr == nil {
		accOnCycleRow := model.DeviceKPI{
			Name:               deviceKpiAccOnCycle21,
			MatchedDeviceCount: cycleAccOnMatched,
			TotalDeviceCount:   cycleTotalVehicles,
			Timestamp:          timestamp,
		}
		if insertErr := db.InsertDeviceKPI(accOnCycleRow); insertErr != nil {
			log.Printf("[ERROR] device acc job: insert ACC_ON_OVER_10H_CYCLE_21 row: %v", insertErr)
			insertFailed = true
		}

		accOffCycleRow := model.DeviceKPI{
			Name:               deviceKpiAccOffCycle21,
			MatchedDeviceCount: cycleAccOffMatched,
			TotalDeviceCount:   cycleTotalVehicles,
			Timestamp:          timestamp,
		}
		if insertErr := db.InsertDeviceKPI(accOffCycleRow); insertErr != nil {
			log.Printf("[ERROR] device acc job: insert ACC_OFF_OVER_72H_CYCLE_21 row: %v", insertErr)
			insertFailed = true
		}

		mknCycleRow := model.DeviceKPI{
			Name:               deviceKpiMknCycle21,
			MatchedDeviceCount: cycleMknMatched,
			TotalDeviceCount:   cycleTotalVehicles,
			Timestamp:          timestamp,
		}
		if insertErr := db.InsertDeviceKPI(mknCycleRow); insertErr != nil {
			log.Printf("[ERROR] device acc job: insert VEHICLE_MKN_OVER_72H_CYCLE_21 row: %v", insertErr)
			insertFailed = true
		}
		log.Printf("[INFO] device acc job: cycle-21 KPI window start_ts=%d total=%d acc_on=%d acc_off=%d mkn=%d",
			cycleStartTs, cycleTotalVehicles, cycleAccOnMatched, cycleAccOffMatched, cycleMknMatched)
	}

	if insertFailed {
		log.Printf("[WARN] device acc job: one or more device_kpi inserts failed; hourly cursor not advanced")
		return
	}
	if err := markDeviceKpiSnapshotInserted(db, timestamp); err != nil {
		log.Printf("[ERROR] device acc job: update hourly insert cursor for device_kpi: %v", err)
	}
}

func countVehicleByFieldCondition(fieldPrefix, fieldValue string, overDuration time.Duration) (int64, int64, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return 0, 0, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)

	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return 0, 0, fmt.Errorf("connect mongo with fallback auth: %w", err)
	}

	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()
	log.Printf("[INFO] device acc job: connected mongo using authSource=%s authMechanism=%s", usedAuthSource, usedAuthMech)

	coll := client.Database(mongoDBName).Collection("attributes")
	threshold := time.Now().UTC().Add(-overDuration).UnixMilli()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	matchedFilter := bson.M{
		"entity_type":                                 "VEHICLE",
		fmt.Sprintf("%s.value", fieldPrefix):          fieldValue,
		fmt.Sprintf("%s.last_update_ts", fieldPrefix): bson.M{"$lte": threshold},
	}

	matchedCount, err := coll.CountDocuments(ctx, matchedFilter)
	if err != nil {
		return 0, 0, fmt.Errorf("count matched devices: %w", err)
	}

	totalFilter := bson.M{
		"entity_type": "VEHICLE",
	}

	totalCount, err := coll.CountDocuments(ctx, totalFilter)
	if err != nil {
		return 0, 0, fmt.Errorf("count total vehicles: %w", err)
	}

	return matchedCount, totalCount, nil
}

func syncDailyVehicleStatusAndGetCounts(db *model.PostgreDb) (int64, int64, int64, int64, error) {
	if err := db.EnsureDeviceKPIDailyStatusTable(); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("ensure device_kpi_daily_vehicle_status table: %w", err)
	}
	if err := db.EnsureDeviceKPIDailyStatusArchiveTable(); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("ensure device_kpi_daily_vehicle_status_archive table: %w", err)
	}
	if err := db.EnsureDeviceKPIJobCursorTable(); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("ensure device_kpi_job_cursor table: %w", err)
	}

	loc := resolveKpiLocation()
	nowLocal := time.Now().In(loc)
	kpiDate := nowLocal.Format("2006-01-02")
	runEndTs := nowLocal.UTC().UnixMilli()
	dayStartTs := startOfDayMillis(nowLocal)

	if err := resetDeviceKPIDailyStatusFlagsOnNewDay(db, nowLocal); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("reset daily status seen flags on new day: %w", err)
	}

	cursor, err := db.GetDeviceKPIJobCursor(deviceDailyStatusCursor)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("get cursor for device daily status job: %w", err)
	}

	lookbackMillis := resolveKpiIncrementalLookbackMillis()
	windowStartTs := dayStartTs
	inclusiveStart := true
	if cursor.Valid && cursor.Int64 >= dayStartTs {
		windowStartTs = cursor.Int64 - lookbackMillis
		if windowStartTs < dayStartTs {
			windowStartTs = dayStartTs
		}
	}
	if windowStartTs > runEndTs {
		windowStartTs = runEndTs
	}

	offlineVehicles, badGpsVehicles, badMemVehicles, onlineVehicles, err := fetchDailyVehicleStatusPlateNoSnapshot(windowStartTs, runEndTs, dayStartTs, inclusiveStart)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, offlineVehicles, "offline"); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("upsert offline plates: %w", err)
	}
	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, badGpsVehicles, "badgps"); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("upsert badgps plates: %w", err)
	}
	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, badMemVehicles, "bad_memory"); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("upsert bad memory plates: %w", err)
	}
	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, onlineVehicles, "online"); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("upsert online plates: %w", err)
	}

	offlineCount, badGpsCount, badMemCount, onlineCount, err := db.CountDeviceKPIDailyStatusFlags(kpiDate)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("count daily status flags: %w", err)
	}

	if err := db.UpsertDeviceKPIJobCursor(deviceDailyStatusCursor, runEndTs); err != nil {
		return 0, 0, 0, 0, fmt.Errorf("upsert cursor for device daily status job: %w", err)
	}

	return offlineCount, badGpsCount, badMemCount, onlineCount, nil
}

func fetchDailyVehicleStatusPlateNoSnapshot(startTs, endTs, dayStartTs int64, inclusiveStart bool) ([]model.DeviceKPIDailyStatusVehicle, []model.DeviceKPIDailyStatusVehicle, []model.DeviceKPIDailyStatusVehicle, []model.DeviceKPIDailyStatusVehicle, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)

	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("connect mongo with fallback auth: %w", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()
	log.Printf("[INFO] device acc job: connected mongo using authSource=%s authMechanism=%s", usedAuthSource, usedAuthMech)

	coll := client.Database(mongoDBName).Collection("attributes")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	timeRange := buildIncrementalTimeRange(startTs, endTs, inclusiveStart)

	offlineFilter := bson.M{
		"entity_type":           "VEHICLE",
		"status.value":          statusOffline,
		"status.last_update_ts": timeRange,
		"plateNo.value":         bson.M{"$exists": true, "$ne": ""},
	}
	badGpsFilter := bson.M{
		"entity_type":           "VEHICLE",
		"status.value":          statusBadGPS,
		"status.last_update_ts": timeRange,
		"plateNo.value":         bson.M{"$exists": true, "$ne": ""},
	}
	onlineFilter := bson.M{
		"entity_type":           "VEHICLE",
		"status.value":          bson.M{"$in": onlineStatusValues},
		"status.last_update_ts": timeRange,
		"plateNo.value":         bson.M{"$exists": true, "$ne": ""},
	}
	badMemoryFilter := bson.M{
		"entity_type":                          "VEHICLE",
		"plateNo.value":                        bson.M{"$exists": true, "$ne": ""},
		"memoryCardStatus.last_update_ts":      timeRange,
		"memoryCardStatus.value.nok_latest_ts": bson.M{"$gte": dayStartTs, "$lte": endTs},
	}

	offlineVehicles, err := distinctVehicleIdentities(ctx, coll, offlineFilter, "status.last_update_ts")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("distinct offline plateNo in day: %w", err)
	}
	badGpsVehicles, err := distinctVehicleIdentities(ctx, coll, badGpsFilter, "status.last_update_ts")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("distinct badgps plateNo in day: %w", err)
	}
	badMemVehicles, err := distinctVehicleIdentities(ctx, coll, badMemoryFilter, "memoryCardStatus.last_update_ts")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("distinct bad memory plateNo in day: %w", err)
	}
	onlineVehicles, err := distinctVehicleIdentities(ctx, coll, onlineFilter, "status.last_update_ts")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("distinct online plateNo in day: %w", err)
	}

	return offlineVehicles, badGpsVehicles, badMemVehicles, onlineVehicles, nil
}

func countDistinctVehicleKpiInCycle21Window() (int64, int64, int64, int64, int64, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	loc := resolveKpiLocation()
	nowLocal := time.Now().In(loc)
	resetDay := resolveKpiCycleResetDay()
	cycleStartTs := startOfKpiCycleMillis(nowLocal, resetDay)
	accOnThresholdTs := nowLocal.UTC().Add(-10 * time.Hour).UnixMilli()
	over72hThresholdTs := nowLocal.UTC().Add(-72 * time.Hour).UnixMilli()

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return 0, 0, 0, 0, cycleStartTs, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)

	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return 0, 0, 0, 0, cycleStartTs, fmt.Errorf("connect mongo with fallback auth: %w", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()
	log.Printf("[INFO] device acc job: connected mongo using authSource=%s authMechanism=%s", usedAuthSource, usedAuthMech)

	coll := client.Database(mongoDBName).Collection("attributes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	totalFilter := bson.M{
		"entity_type":   "VEHICLE",
		"plateNo.value": bson.M{"$exists": true, "$ne": ""},
	}
	totalVehicles, err := countDistinctPlateNos(ctx, coll, totalFilter)
	if err != nil {
		return 0, 0, 0, 0, cycleStartTs, fmt.Errorf("distinct total vehicles in cycle window: %w", err)
	}

	var accOnCount int64
	if accOnThresholdTs >= cycleStartTs {
		accOnFilter := buildDistinctPlateWindowFilter("acc", "active", cycleStartTs, accOnThresholdTs)
		accOnCount, err = countDistinctPlateNos(ctx, coll, accOnFilter)
		if err != nil {
			return 0, 0, 0, 0, cycleStartTs, fmt.Errorf("distinct ACC_ON_OVER_10H in cycle window: %w", err)
		}
	}

	var accOffCount int64
	if over72hThresholdTs >= cycleStartTs {
		accOffFilter := buildDistinctPlateWindowFilter("acc", "inactive", cycleStartTs, over72hThresholdTs)
		accOffCount, err = countDistinctPlateNos(ctx, coll, accOffFilter)
		if err != nil {
			return 0, 0, 0, 0, cycleStartTs, fmt.Errorf("distinct ACC_OFF_OVER_72H in cycle window: %w", err)
		}
	}

	var mknCount int64
	if over72hThresholdTs >= cycleStartTs {
		mknFilter := buildDistinctPlateWindowFilter("status", statusOffline, cycleStartTs, over72hThresholdTs)
		mknCount, err = countDistinctPlateNos(ctx, coll, mknFilter)
		if err != nil {
			return 0, 0, 0, 0, cycleStartTs, fmt.Errorf("distinct VEHICLE_MKN_OVER_72H in cycle window: %w", err)
		}
	}

	return accOnCount, accOffCount, mknCount, totalVehicles, cycleStartTs, nil
}

type mongoVehicleDetailRow struct {
	PlateNo      string `bson:"plate_no"`
	IMEI         string `bson:"imei"`
	LastUpdateTS int64  `bson:"last_update_ts"`
}

type mongoVehicleACCStateRow struct {
	PlateNo         string `bson:"plate_no"`
	IMEI            string `bson:"imei"`
	ACCValue        string `bson:"acc_value"`
	ACCLastUpdateTS int64  `bson:"acc_last_update_ts"`
}

type mongoVehicleIdentityRow struct {
	PlateNo string `bson:"plate_no"`
	IMEI    string `bson:"imei"`
	LastTS  int64  `bson:"last_ts"`
}

func fetchVehicleACCStateSnapshot(sampledAtUTC time.Time) ([]model.DeviceVehicleACCState, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return nil, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)
	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return nil, fmt.Errorf("connect mongo with fallback auth: %w", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()
	log.Printf("[INFO] device acc job: connected mongo using authSource=%s authMechanism=%s", usedAuthSource, usedAuthMech)

	coll := client.Database(mongoDBName).Collection("attributes")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"entity_type":   "VEHICLE",
			"plateNo.value": bson.M{"$exists": true, "$ne": ""},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":                0,
			"plate_no":           "$plateNo.value",
			"imei":               bson.M{"$ifNull": []interface{}{"$imei.value", ""}},
			"acc_value":          bson.M{"$ifNull": []interface{}{"$acc.value", ""}},
			"acc_last_update_ts": bson.M{"$ifNull": []interface{}{"$acc.last_update_ts", 0}},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate ACC state snapshot from mongo: %w", err)
	}
	defer cursor.Close(ctx)

	sampledAtUTC = sampledAtUTC.UTC()
	byPlate := make(map[string]model.DeviceVehicleACCState)
	for cursor.Next(ctx) {
		var doc mongoVehicleACCStateRow
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode ACC state row: %w", err)
		}

		plate := strings.TrimSpace(doc.PlateNo)
		if plate == "" {
			continue
		}
		ts := doc.ACCLastUpdateTS
		if ts <= 0 {
			ts = sampledAtUTC.UnixMilli()
		}
		accValue := strings.ToLower(strings.TrimSpace(doc.ACCValue))
		if accValue != "active" && accValue != "inactive" {
			accValue = ""
		}

		row := model.DeviceVehicleACCState{
			PlateNo:      plate,
			IMEI:         strings.TrimSpace(doc.IMEI),
			ACCValue:     accValue,
			LastUpdateTS: ts,
			LastUpdateAt: time.UnixMilli(ts).UTC(),
		}
		if existing, ok := byPlate[plate]; !ok || row.LastUpdateTS > existing.LastUpdateTS {
			byPlate[plate] = row
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate ACC state rows: %w", err)
	}

	rows := make([]model.DeviceVehicleACCState, 0, len(byPlate))
	for _, row := range byPlate {
		rows = append(rows, row)
	}
	return rows, nil
}

func fetchVehicleKpiDetailSnapshot(sampledAtUTC time.Time) ([]model.DeviceKPIVehicleDetail, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return nil, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)
	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return nil, fmt.Errorf("connect mongo with fallback auth: %w", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()
	log.Printf("[INFO] device acc job: connected mongo using authSource=%s authMechanism=%s", usedAuthSource, usedAuthMech)

	sampledAtUTC = sampledAtUTC.UTC()
	sampledHourUTC := sampledHourFromUTC(sampledAtUTC, resolveKpiLocation())
	over72hThresholdTs := sampledAtUTC.Add(-72 * time.Hour).UnixMilli()

	coll := client.Database(mongoDBName).Collection("attributes")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	mknRows, err := fetchVehicleDetailsByFieldCondition(ctx, coll, deviceKpiVehicleMknName, "status", statusOffline, over72hThresholdTs, sampledAtUTC, sampledHourUTC)
	if err != nil {
		return nil, err
	}
	return mknRows, nil
}

func sampledHourFromUTC(sampledAtUTC time.Time, loc *time.Location) time.Time {
	local := sampledAtUTC.UTC().In(loc)
	localHour := time.Date(local.Year(), local.Month(), local.Day(), local.Hour(), 0, 0, 0, loc)
	return localHour.UTC()
}

func resetDeviceKPIVehicleDetailFlagsOnNewDay(db *model.PostgreDb, nowUTC time.Time) error {
	if err := db.EnsureDeviceKPIJobCursorTable(); err != nil {
		return fmt.Errorf("ensure cursor table: %w", err)
	}

	loc := resolveKpiLocation()
	nowLocal := nowUTC.UTC().In(loc)
	dayStartLocal := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc)
	dayStartUTC := dayStartLocal.UTC()
	previousSnapshotDate := dayStartLocal.AddDate(0, 0, -1).Format("2006-01-02")

	cursor, err := db.GetDeviceKPIJobCursor(deviceKpiDetailDayReset)
	if err != nil {
		return fmt.Errorf("get day reset cursor: %w", err)
	}
	if cursor.Valid && cursor.Int64 >= dayStartUTC.UnixMilli() {
		return nil
	}

	if err := db.SnapshotDeviceKPIVehicleDetail(previousSnapshotDate, nowUTC); err != nil {
		return fmt.Errorf("snapshot device_kpi_vehicle_detail for snapshot_date=%s: %w", previousSnapshotDate, err)
	}
	sampledHourUTC := sampledHourFromUTC(nowUTC, loc)
	if err := db.ResetAllDeviceKPIVehicleDetailFlags(nowUTC, sampledHourUTC); err != nil {
		return fmt.Errorf("reset all detail flags: %w", err)
	}
	if err := db.UpsertDeviceKPIJobCursor(deviceKpiDetailDayReset, dayStartUTC.UnixMilli()); err != nil {
		return fmt.Errorf("upsert day reset cursor: %w", err)
	}

	log.Printf("[INFO] device acc job: snapshot device_kpi_vehicle_detail snapshot_date=%s and reset flags for new day=%s",
		previousSnapshotDate, dayStartLocal.Format("2006-01-02"))
	return nil
}

func resetDeviceKPIDailyStatusFlagsOnNewDay(db *model.PostgreDb, nowLocal time.Time) error {
	loc := nowLocal.Location()
	dayStartLocal := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc)
	dayStartUTC := dayStartLocal.UTC()
	kpiDate := dayStartLocal.Format("2006-01-02")
	previousKpiDate := dayStartLocal.AddDate(0, 0, -1).Format("2006-01-02")

	cursor, err := db.GetDeviceKPIJobCursor(deviceDailyStatusReset)
	if err != nil {
		return fmt.Errorf("get day reset cursor: %w", err)
	}
	if cursor.Valid && cursor.Int64 >= dayStartUTC.UnixMilli() {
		return nil
	}

	if err := db.SnapshotDeviceKPIDailyStatus(previousKpiDate, nowLocal.UTC()); err != nil {
		return fmt.Errorf("snapshot daily status for kpi_date=%s: %w", previousKpiDate, err)
	}
	deletedRows, err := db.DeleteDeviceKPIDailyStatusByDate(previousKpiDate)
	if err != nil {
		return fmt.Errorf("delete runtime daily status rows for kpi_date=%s: %w", previousKpiDate, err)
	}
	if err := db.ResetDeviceKPIDailyStatusFlags(kpiDate); err != nil {
		return fmt.Errorf("reset seen flags for kpi_date=%s: %w", kpiDate, err)
	}
	if err := db.UpsertDeviceKPIJobCursor(deviceDailyStatusReset, dayStartUTC.UnixMilli()); err != nil {
		return fmt.Errorf("upsert day reset cursor: %w", err)
	}

	log.Printf("[INFO] device acc job: snapshot+cleanup device_kpi_daily_vehicle_status kpi_date=%s deleted_rows=%d and reset flags for new day=%s", previousKpiDate, deletedRows, kpiDate)
	return nil
}

func fetchVehicleDetailsByFieldCondition(
	ctx context.Context,
	coll *mongo.Collection,
	kpiName, fieldPrefix, fieldValue string,
	thresholdTs int64,
	sampledAtUTC, sampledHourUTC time.Time,
) ([]model.DeviceKPIVehicleDetail, error) {
	filter := bson.M{
		"entity_type":                                 "VEHICLE",
		"plateNo.value":                               bson.M{"$exists": true, "$ne": ""},
		fmt.Sprintf("%s.value", fieldPrefix):          fieldValue,
		fmt.Sprintf("%s.last_update_ts", fieldPrefix): bson.M{"$lte": thresholdTs},
	}

	project := bson.M{
		"_id":            0,
		"plate_no":       "$plateNo.value",
		"imei":           bson.M{"$ifNull": []interface{}{"$imei.value", ""}},
		"last_update_ts": fmt.Sprintf("$%s.last_update_ts", fieldPrefix),
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$project", Value: project}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate %s from mongo: %w", kpiName, err)
	}
	defer cursor.Close(ctx)

	plateMap := make(map[string]model.DeviceKPIVehicleDetail)
	for cursor.Next(ctx) {
		var doc mongoVehicleDetailRow
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode %s vehicle detail row: %w", kpiName, err)
		}

		plateNo := strings.TrimSpace(doc.PlateNo)
		if plateNo == "" || doc.LastUpdateTS <= 0 {
			continue
		}

		imei := strings.TrimSpace(doc.IMEI)
		row := model.DeviceKPIVehicleDetail{
			KPIName:      kpiName,
			PlateNo:      plateNo,
			IMEI:         imei,
			LastUpdateTS: doc.LastUpdateTS,
			LastUpdateAt: time.UnixMilli(doc.LastUpdateTS).UTC(),
			SampledAt:    sampledAtUTC,
			SampledHour:  sampledHourUTC,
		}

		if existing, ok := plateMap[plateNo]; !ok || row.LastUpdateTS > existing.LastUpdateTS {
			plateMap[plateNo] = row
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s vehicle detail rows: %w", kpiName, err)
	}

	rows := make([]model.DeviceKPIVehicleDetail, 0, len(plateMap))
	for _, row := range plateMap {
		rows = append(rows, row)
	}

	return rows, nil
}

func buildDistinctPlateWindowFilter(fieldPrefix, fieldValue string, startTs, endTs int64) bson.M {
	return bson.M{
		"entity_type":                                 "VEHICLE",
		"plateNo.value":                               bson.M{"$exists": true, "$ne": ""},
		fmt.Sprintf("%s.value", fieldPrefix):          fieldValue,
		fmt.Sprintf("%s.last_update_ts", fieldPrefix): bson.M{"$gte": startTs, "$lte": endTs},
	}
}

func distinctPlateNos(ctx context.Context, coll *mongo.Collection, filter bson.M) ([]string, error) {
	values, err := coll.Distinct(ctx, "plateNo.value", filter)
	if err != nil {
		return nil, err
	}

	plateNos := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, raw := range values {
		var plate string
		if s, ok := raw.(string); ok {
			plate = strings.TrimSpace(s)
		} else {
			plate = strings.TrimSpace(fmt.Sprintf("%v", raw))
		}
		if plate == "" || plate == "<nil>" {
			continue
		}
		if _, exists := seen[plate]; exists {
			continue
		}
		seen[plate] = struct{}{}
		plateNos = append(plateNos, plate)
	}

	return plateNos, nil
}

func countDistinctPlateNos(ctx context.Context, coll *mongo.Collection, filter bson.M) (int64, error) {
	plateNos, err := distinctPlateNos(ctx, coll, filter)
	if err != nil {
		return 0, err
	}

	return int64(len(plateNos)), nil
}

func distinctVehicleIdentities(ctx context.Context, coll *mongo.Collection, filter bson.M, tsField string) ([]model.DeviceKPIDailyStatusVehicle, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$project", Value: bson.M{
			"_id":      0,
			"plate_no": "$plateNo.value",
			"imei":     bson.M{"$ifNull": []interface{}{"$imei.value", ""}},
			"last_ts":  fmt.Sprintf("$%s", tsField),
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "plate_no", Value: 1},
			{Key: "last_ts", Value: -1},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":     "$plate_no",
			"imei":    bson.M{"$first": "$imei"},
			"last_ts": bson.M{"$first": "$last_ts"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":      0,
			"plate_no": "$_id",
			"imei":     1,
			"last_ts":  1,
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	vehicles := make([]model.DeviceKPIDailyStatusVehicle, 0)
	for cursor.Next(ctx) {
		var row mongoVehicleIdentityRow
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		plate := strings.TrimSpace(row.PlateNo)
		if plate == "" {
			continue
		}
		vehicles = append(vehicles, model.DeviceKPIDailyStatusVehicle{
			PlateNo: plate,
			IMEI:    strings.TrimSpace(row.IMEI),
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func startOfDayMillis(nowLocal time.Time) int64 {
	startOfDayLocal := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, nowLocal.Location())
	return startOfDayLocal.UTC().UnixMilli()
}

func startOfKpiCycleMillis(nowLocal time.Time, resetDay int) int64 {
	year, month, _ := nowLocal.Date()
	if nowLocal.Day() <= resetDay {
		month--
	}
	start := time.Date(year, month, resetDay, 0, 0, 0, 0, nowLocal.Location())
	return start.UTC().UnixMilli()
}

func resolveKpiLocation() *time.Location {
	tz := strings.TrimSpace(config.Env(config.EnvKpiTimezone, config.DefKpiTimezone))
	if tz == "" {
		tz = config.DefKpiTimezone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("[WARN] device acc job: invalid KPI timezone=%q, fallback=%q", tz, config.DefKpiTimezone)
		fallback, fbErr := time.LoadLocation(config.DefKpiTimezone)
		if fbErr != nil {
			log.Printf("[WARN] device acc job: cannot load default timezone=%q, fallback to fixed UTC+7", config.DefKpiTimezone)
			return time.FixedZone("UTC+7", 7*60*60)
		}
		return fallback
	}

	return loc
}

func resolveKpiIncrementalLookbackMillis() int64 {
	raw := strings.TrimSpace(config.Env(config.EnvKpiIncrementalLookbackSecond, config.DefKpiIncrementalLookbackSecond))
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || seconds < 0 {
		log.Printf("[WARN] device acc job: invalid KPI incremental lookback seconds=%q, fallback=%s", raw, config.DefKpiIncrementalLookbackSecond)
		fallback, fbErr := strconv.ParseInt(config.DefKpiIncrementalLookbackSecond, 10, 64)
		if fbErr != nil || fallback < 0 {
			return 0
		}
		return fallback * 1000
	}

	return seconds * 1000
}

func resolveKpiCycleResetDay() int {
	raw := strings.TrimSpace(config.Env(config.EnvKpiCycleResetDay, config.DefKpiCycleResetDay))
	day, err := strconv.Atoi(raw)
	if err != nil || day < 1 || day > 28 {
		log.Printf("[WARN] device acc job: invalid KPI cycle reset day=%q, fallback=%s", raw, config.DefKpiCycleResetDay)
		fallback, fbErr := strconv.Atoi(config.DefKpiCycleResetDay)
		if fbErr != nil || fallback < 1 || fallback > 28 {
			return 21
		}
		return fallback
	}

	return day
}

func shouldInsertDeviceKpiSnapshot(db *model.PostgreDb, nowUTC time.Time) (bool, time.Time, error) {
	if err := db.EnsureDeviceKPIJobCursorTable(); err != nil {
		return false, time.Time{}, err
	}

	cursor, err := db.GetDeviceKPIJobCursor(deviceKpiInsertCursor)
	if err != nil {
		return false, time.Time{}, err
	}

	if !cursor.Valid {
		return true, nowUTC, nil
	}

	lastInsertedAt := time.UnixMilli(cursor.Int64).UTC()
	nextAllowed := lastInsertedAt.Add(time.Hour)
	if !nowUTC.Before(nextAllowed) {
		return true, nextAllowed, nil
	}

	return false, nextAllowed, nil
}

func markDeviceKpiSnapshotInserted(db *model.PostgreDb, insertedAtUTC time.Time) error {
	if err := db.EnsureDeviceKPIJobCursorTable(); err != nil {
		return err
	}
	return db.UpsertDeviceKPIJobCursor(deviceKpiInsertCursor, insertedAtUTC.UnixMilli())
}

func buildIncrementalTimeRange(startTs, endTs int64, inclusiveStart bool) bson.M {
	if inclusiveStart {
		return bson.M{
			"$gte": startTs,
			"$lte": endTs,
		}
	}

	return bson.M{
		"$gt":  startTs,
		"$lte": endTs,
	}
}

func connectMongoWithFallback(uri, user, pass string, authSources, authMechs []string) (*mongo.Client, string, string, error) {
	var lastErr error

	for _, authSource := range authSources {
		for _, authMechanism := range authMechs {
			client, err := connectMongoOnce(uri, user, pass, authSource, authMechanism)
			if err == nil {
				if authMechanism == "" {
					return client, authSource, "auto", nil
				}
				return client, authSource, authMechanism, nil
			}
			lastErr = err
			if authMechanism == "" {
				log.Printf("[WARN] device acc job: mongo connect failed with authSource=%s authMechanism=auto: %v", authSource, err)
			} else {
				log.Printf("[WARN] device acc job: mongo connect failed with authSource=%s authMechanism=%s: %v", authSource, authMechanism, err)
			}
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("no mongo auth candidates available")
	}
	return nil, "", "", lastErr
}

func connectMongoOnce(uri, user, pass, authSource, authMechanism string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credential := options.Credential{
		Username:   user,
		Password:   pass,
		AuthSource: authSource,
	}
	if authMechanism != "" {
		credential.AuthMechanism = authMechanism
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(credential))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
		return nil, err
	}

	return client, nil
}

func buildAuthSourceCandidates(dbName, configured string) []string {
	candidates := make([]string, 0, 3)
	pushUnique := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		if !containsString(candidates, v) {
			candidates = append(candidates, v)
		}
	}

	pushUnique(configured)
	pushUnique(dbName)
	pushUnique("admin")

	return candidates
}

func buildAuthMechanismCandidates(configured string) []string {
	candidates := make([]string, 0, 3)
	pushUnique := func(v string) {
		v = strings.TrimSpace(v)
		if containsString(candidates, v) {
			return
		}
		candidates = append(candidates, v)
	}

	// empty mechanism means letting the driver negotiate automatically.
	pushUnique(configured)
	pushUnique("")
	pushUnique("SCRAM-SHA-256")
	pushUnique("SCRAM-SHA-1")

	return candidates
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func buildMongoURI(hosts, defaultPort string) (string, error) {
	rawHosts := strings.Split(hosts, ",")
	normalizedHosts := make([]string, 0, len(rawHosts))

	for _, host := range rawHosts {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		if strings.Contains(host, ":") {
			normalizedHosts = append(normalizedHosts, host)
			continue
		}
		normalizedHosts = append(normalizedHosts, host+":"+defaultPort)
	}

	if len(normalizedHosts) == 0 {
		return "", fmt.Errorf("mongo hosts are empty")
	}

	return "mongodb://" + strings.Join(normalizedHosts, ","), nil
}
