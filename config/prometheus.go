package config

const PrometheusUrl = "http://prometheus-stack-kube-prom-prometheus.monitoring:9090/api/v1/query"

type KpiData struct {
	Pod    string `json:"pod"`
	Method string `json:"method"`
	Req    string `json:"request"`
	Error  string `json:"error"`
}

type LatencyKpi struct {
	Api        string `json:"api"`
	Total      string `json:"total"`
	Count      string `json:"count"`      //so luong ban ghi nho hon 5s
	Percentile string `json:"percentile"` //" percentile 95 theo yeu cau
	Result     string `json:"result"`
}

type SuccessKpi struct {
	Api   string  `json:"api"`
	Value float64 `json:"value"`
}

type Ramusage struct {
	Server string `json:"server"`
	Total  string `json:"total"`
	Avail  string `json:"avail"`
}

type Cpuusage struct {
	Server string `json:"server"`
	Value  string `json:"value"`
}
