package validations

type GetPrometheusVtrack struct {
	query string `json:"query" binding:"required"`
}
