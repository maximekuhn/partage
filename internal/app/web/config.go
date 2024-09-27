package web

type ServerConfig struct {
	// DBFilepath is the file path to use for SQLite.
	// The file may not exist (it will be created) but must be a valid path.
	DBFilepath string

	// JWT Signature Key (HMAC)
	JWTSignatureKey []byte
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		DBFilepath:      "/tmp/partage.sqlite3",
		JWTSignatureKey: []byte("64a6988ec0ecacbdf40ecf504e70b9a5f6174a8992c856c7ee22e1e0be03a8890412904b9d17a467d03559fe573c324271615dbcf191e4cfc259b5a01a3bb824"),
	}
}
