package config

const (
	EnvKpiTimezone                  = "KPI_TIMEZONE"
	DefKpiTimezone                  = "Asia/Bangkok"
	EnvKpiIncrementalLookbackSecond = "KPI_INCREMENTAL_LOOKBACK_SECONDS"
	DefKpiIncrementalLookbackSecond = "300"
	EnvKpiCycleResetDay             = "KPI_CYCLE_RESET_DAY"
	DefKpiCycleResetDay             = "21"
	EnvKpiDetailUpsertTimeoutSecond = "KPI_DETAIL_UPSERT_TIMEOUT_SECONDS"
	DefKpiDetailUpsertTimeoutSecond = "180"
	EnvKpiDetailUpsertBatchSize     = "KPI_DETAIL_UPSERT_BATCH_SIZE"
	DefKpiDetailUpsertBatchSize     = "5000"
)
