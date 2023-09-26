package validations

type GetMethod struct {
	oid string `json:"oid" binding:"required"`
}
