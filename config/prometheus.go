package config

const PrometheusUrl = "http://kube-prometheus-stack-prometheus.monitoring:9090/api/v1/query"

type KpiData struct {
	Pod    string `json:"pod"`
	Method string `json:"method"`
	Req    string `json:"request"`
	Error  string `json:"error"`
}
