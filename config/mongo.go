package config

const (
	EnvAttrMongoDBName   = "ATTR_MONGO_DB_NAME"
	EnvAttrMongoHosts    = "ATTR_MONGO_HOSTS"
	EnvAttrMongoPort     = "ATTR_MONGO_PORT"
	EnvAttrMongoUser     = "ATTR_MONGO_USER"
	EnvAttrMongoPass     = "ATTR_MONGO_PASS"
	EnvAttrMongoAuthDB   = "ATTR_MONGO_AUTH_SOURCE"
	EnvAttrMongoAuthMech = "ATTR_MONGO_AUTH_MECHANISM"

	DefAttrMongoDBName   = "attributes"
	DefAttrMongoHosts    = "10.207.191.94:27017,10.207.191.95:27017,10.207.191.96:27017"
	DefAttrMongoPort     = "27017"
	DefAttrMongoUser     = "viot"
	DefAttrMongoPass     = "newpassword"
	DefAttrMongoAuthDB   = ""
	DefAttrMongoAuthMech = ""
)
