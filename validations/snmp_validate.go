package validations

type GetOidSnmp struct {
	oid string `json:"oid" binding:"required"`
}

type GetEachOidWithServer struct {
	oid    string `json:"oid" binding:"required"`
	server string `json:"server" binding:"required"`
}
