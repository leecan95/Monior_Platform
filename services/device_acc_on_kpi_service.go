package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"context"
	"fmt"
	"log"
	"runtime/debug"
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
	log.Printf("[INFO] device acc job: worker bootstrapping (interval=5m, insert_interval=1h)")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] device acc job: worker goroutine crashed: %v", r)
				log.Printf("[PANIC] device acc job: stacktrace:\n%s", string(debug.Stack()))
			}
		}()

		runDeviceAccOnKpiSnapshotSafe("startup")

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		log.Printf("[INFO] device acc job: ticker started interval=%s", (5 * time.Minute).String())

		for tickAt := range ticker.C {
			runDeviceAccOnKpiSnapshotSafe("ticker@" + tickAt.UTC().Format(time.RFC3339))
		}
	}()
}

func runDeviceAccOnKpiSnapshotSafe(trigger string) {
	startedAt := time.Now().UTC()
	log.Printf("[INFO] device acc job: snapshot begin trigger=%s started_at=%s", trigger, startedAt.Format(time.RFC3339))
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] device acc job: snapshot panic trigger=%s err=%v", trigger, r)
			log.Printf("[PANIC] device acc job: stacktrace:\n%s", string(debug.Stack()))
		} else {
			log.Printf("[INFO] device acc job: snapshot end trigger=%s duration=%s", trigger, time.Since(startedAt))
		}
	}()

	runDeviceAccOnKpiSnapshot()
}

func runDeviceAccOnKpiSnapshot() {
	cfg := model.LoadDBConfig()
	log.Printf("[INFO] device acc job: loaded DB config host=%s port=%s user=%s db=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name)

	matchedAccOn, totalAccOn, err := countVehicleByFieldCondition("acc", "active", 10*time.Hour)
	if err != nil {
		log.Printf("[ERROR] device acc job: count ACC_ON_OVER_10H from mongo: %v", err)
		return
	}
	matchedAccOff, totalAccOff, err := countVehicleByFieldCondition("acc", "inactive", 72*time.Hour)
	if err != nil {
		log.Printf("[ERROR] device acc job: count ACC_OFF_OVER_72H from mongo: %v", err)
		return
	}
	matchedMkn, totalMkn, err := countVehicleByFieldCondition("status", "offline", 72*time.Hour)
	if err != nil {
		log.Printf("[ERROR] device acc job: count VEHICLE_MKN_OVER_72H from mongo: %v", err)
		return
	}
	log.Printf("[INFO] device acc job: mongo counters acc_on=%d/%d acc_off=%d/%d mkn_72h=%d/%d",
		matchedAccOn, totalAccOn, matchedAccOff, totalAccOff, matchedMkn, totalMkn)

	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] device acc job: connect trans db: %v", err)
		return
	}
	log.Printf("[INFO] device acc job: connected transactions DB")
	if err := db.EnsureDeviceKPITable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi table: %v", err)
		return
	}
	log.Printf("[INFO] device acc job: ensured table device_kpi")

	matchedMknCurrent, matchedBadGps, totalOnlineInDay, err := syncDailyVehicleStatusAndGetCounts(db)
	if err != nil {
		log.Printf("[ERROR] device acc job: sync daily status for VEHICLE_MKN_CURRENT and VEHICLE_BADGPS_OVER_ONLINE_IN_DAY: %v", err)
		return
	}
	log.Printf("[INFO] device acc job: daily counters mkn_current=%d badgps_in_day=%d online_in_day=%d",
		matchedMknCurrent, matchedBadGps, totalOnlineInDay)

	cycleAccOnMatched, cycleAccOffMatched, cycleMknMatched, cycleTotalVehicles, cycleStartTs, cycleErr := countDistinctVehicleKpiInCycle21Window()
	if cycleErr != nil {
		log.Printf("[ERROR] device acc job: count cycle-21 KPI from mongo: %v", cycleErr)
	}

	timestamp := time.Now().UTC()
	canInsertDeviceKpi, insertWait, err := shouldInsertDeviceKpiSnapshot(db, timestamp)
	if err != nil {
		log.Printf("[ERROR] device acc job: decide hourly insert for device_kpi: %v", err)
		return
	}
	if !canInsertDeviceKpi {
		log.Printf("[INFO] device acc job: sampled mongo, skip device_kpi insert (next write after %s)", insertWait.UTC().Format(time.RFC3339))
		return
	}
	log.Printf("[INFO] device acc job: hourly gate passed, proceed writing snapshots")

	if err := db.EnsureDeviceKPIVehicleDetailTable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi_vehicle_detail table: %v", err)
		return
	}
	if err := db.EnsureDeviceKPIVehicleDetailSnapshotTable(); err != nil {
		log.Printf("[ERROR] device acc job: ensure device_kpi_vehicle_detail_snapshot table: %v", err)
		return
	}
	if err := resetDeviceKPIVehicleDetailFlagsOnNewDay(db, timestamp); err != nil {
		log.Printf("[ERROR] device acc job: reset daily flags in device_kpi_vehicle_detail: %v", err)
		return
	}
	detailRows, err := fetchVehicleKpiDetailSnapshot(timestamp)
	if err != nil {
		log.Printf("[ERROR] device acc job: fetch hourly vehicle detail snapshot from mongo: %v", err)
		return
	}
	sampledHourUTC := sampledHourFromUTC(timestamp, resolveKpiLocation())
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
		return
	}
	log.Printf("[INFO] device acc job: device_kpi hourly snapshot committed at=%s", timestamp.Format(time.RFC3339))
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
	log.Printf("[INFO] device acc job: count condition field=%s value=%s threshold_ts=%d", fieldPrefix, fieldValue, threshold)
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

func syncDailyVehicleStatusAndGetCounts(db *model.PostgreDb) (int64, int64, int64, error) {
	if err := db.EnsureDeviceKPIDailyStatusTable(); err != nil {
		return 0, 0, 0, fmt.Errorf("ensure device_kpi_daily_vehicle_status table: %w", err)
	}
	if err := db.EnsureDeviceKPIDailyStatusArchiveTable(); err != nil {
		return 0, 0, 0, fmt.Errorf("ensure device_kpi_daily_vehicle_status_archive table: %w", err)
	}
	if err := db.EnsureDeviceKPIJobCursorTable(); err != nil {
		return 0, 0, 0, fmt.Errorf("ensure device_kpi_job_cursor table: %w", err)
	}

	loc := resolveKpiLocation()
	nowLocal := time.Now().In(loc)
	kpiDate := nowLocal.Format("2006-01-02")
	runEndTs := nowLocal.UTC().UnixMilli()
	dayStartTs := startOfDayMillis(nowLocal)

	if err := resetDeviceKPIDailyStatusFlagsOnNewDay(db, nowLocal); err != nil {
		return 0, 0, 0, fmt.Errorf("reset daily status seen flags on new day: %w", err)
	}

	cursor, err := db.GetDeviceKPIJobCursor(deviceDailyStatusCursor)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get cursor for device daily status job: %w", err)
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
	log.Printf("[INFO] device acc job: daily window kpi_date=%s start_ts=%d end_ts=%d lookback_ms=%d cursor_valid=%t cursor_ts=%d",
		kpiDate, windowStartTs, runEndTs, lookbackMillis, cursor.Valid, cursor.Int64)

	offlinePlates, badGpsPlates, onlinePlates, err := fetchDailyVehicleStatusPlateNoSnapshot(windowStartTs, runEndTs, inclusiveStart)
	if err != nil {
		return 0, 0, 0, err
	}

	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, offlinePlates, "offline"); err != nil {
		return 0, 0, 0, fmt.Errorf("upsert offline plates: %w", err)
	}
	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, badGpsPlates, "badgps"); err != nil {
		return 0, 0, 0, fmt.Errorf("upsert badgps plates: %w", err)
	}
	if err := db.UpsertDeviceKPIDailyStatusFlags(kpiDate, onlinePlates, "online"); err != nil {
		return 0, 0, 0, fmt.Errorf("upsert online plates: %w", err)
	}

	offlineCount, badGpsCount, onlineCount, err := db.CountDeviceKPIDailyStatusFlags(kpiDate)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("count daily status flags: %w", err)
	}

	if err := db.UpsertDeviceKPIJobCursor(deviceDailyStatusCursor, runEndTs); err != nil {
		return 0, 0, 0, fmt.Errorf("upsert cursor for device daily status job: %w", err)
	}
	log.Printf("[INFO] device acc job: daily status cursor advanced to=%d", runEndTs)

	return offlineCount, badGpsCount, onlineCount, nil
}

func fetchDailyVehicleStatusPlateNoSnapshot(startTs, endTs int64, inclusiveStart bool) ([]string, []string, []string, error) {
	mongoDBName := config.Env(config.EnvAttrMongoDBName, config.DefAttrMongoDBName)
	mongoHosts := config.Env(config.EnvAttrMongoHosts, config.DefAttrMongoHosts)
	mongoPort := config.Env(config.EnvAttrMongoPort, config.DefAttrMongoPort)
	mongoUser := config.Env(config.EnvAttrMongoUser, config.DefAttrMongoUser)
	mongoPass := config.Env(config.EnvAttrMongoPass, config.DefAttrMongoPass)
	mongoAuthSourceCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthDB, config.DefAttrMongoAuthDB))
	mongoAuthMechCfg := strings.TrimSpace(config.Env(config.EnvAttrMongoAuthMech, config.DefAttrMongoAuthMech))

	clientURI, err := buildMongoURI(mongoHosts, mongoPort)
	if err != nil {
		return nil, nil, nil, err
	}

	authSources := buildAuthSourceCandidates(mongoDBName, mongoAuthSourceCfg)
	authMechs := buildAuthMechanismCandidates(mongoAuthMechCfg)

	client, usedAuthSource, usedAuthMech, err := connectMongoWithFallback(clientURI, mongoUser, mongoPass, authSources, authMechs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connect mongo with fallback auth: %w", err)
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

	offlinePlates, err := distinctPlateNos(ctx, coll, offlineFilter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("distinct offline plateNo in day: %w", err)
	}
	badGpsPlates, err := distinctPlateNos(ctx, coll, badGpsFilter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("distinct badgps plateNo in day: %w", err)
	}
	onlinePlates, err := distinctPlateNos(ctx, coll, onlineFilter)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("distinct online plateNo in day: %w", err)
	}
	log.Printf("[INFO] device acc job: daily distinct plate counts offline=%d badgps=%d online=%d", len(offlinePlates), len(badGpsPlates), len(onlinePlates))

	return offlinePlates, badGpsPlates, onlinePlates, nil
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
	accOnThresholdTs := sampledAtUTC.Add(-10 * time.Hour).UnixMilli()
	over72hThresholdTs := sampledAtUTC.Add(-72 * time.Hour).UnixMilli()

	coll := client.Database(mongoDBName).Collection("attributes")
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	accOnRows, err := fetchVehicleDetailsByFieldCondition(ctx, coll, deviceKpiAccOnName, "acc", "active", accOnThresholdTs, sampledAtUTC, sampledHourUTC)
	if err != nil {
		return nil, err
	}
	accOffRows, err := fetchVehicleDetailsByFieldCondition(ctx, coll, deviceKpiAccOffName, "acc", "inactive", over72hThresholdTs, sampledAtUTC, sampledHourUTC)
	if err != nil {
		return nil, err
	}
	mknRows, err := fetchVehicleDetailsByFieldCondition(ctx, coll, deviceKpiVehicleMknName, "status", statusOffline, over72hThresholdTs, sampledAtUTC, sampledHourUTC)
	if err != nil {
		return nil, err
	}

	rows := make([]model.DeviceKPIVehicleDetail, 0, len(accOnRows)+len(accOffRows)+len(mknRows))
	rows = append(rows, accOnRows...)
	rows = append(rows, accOffRows...)
	rows = append(rows, mknRows...)

	return rows, nil
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
		log.Printf("[INFO] device acc job: skip detail day reset, cursor already at day_start=%d", dayStartUTC.UnixMilli())
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
	log.Printf("[INFO] device acc job: updated cursor %s=%d", deviceKpiDetailDayReset, dayStartUTC.UnixMilli())

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
		log.Printf("[INFO] device acc job: skip daily_status day reset, cursor already at day_start=%d", dayStartUTC.UnixMilli())
		return nil
	}

	if err := db.SnapshotDeviceKPIDailyStatus(previousKpiDate, nowLocal.UTC()); err != nil {
		return fmt.Errorf("snapshot daily status for kpi_date=%s: %w", previousKpiDate, err)
	}
	if err := db.ResetDeviceKPIDailyStatusFlags(kpiDate); err != nil {
		return fmt.Errorf("reset seen flags for kpi_date=%s: %w", kpiDate, err)
	}
	if err := db.UpsertDeviceKPIJobCursor(deviceDailyStatusReset, dayStartUTC.UnixMilli()); err != nil {
		return fmt.Errorf("upsert day reset cursor: %w", err)
	}
	log.Printf("[INFO] device acc job: updated cursor %s=%d", deviceDailyStatusReset, dayStartUTC.UnixMilli())

	log.Printf("[INFO] device acc job: snapshot device_kpi_daily_vehicle_status kpi_date=%s and reset flags for new day=%s", previousKpiDate, kpiDate)
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
	candidates := make([]string, 0, 4)
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
