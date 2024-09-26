package web

type ServerConfig struct {
	// DBFilepath is the file path to use for SQLite.
	// The file may not exist (it will be created) but must be a valid path.
	DBFilepath string
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		DBFilepath: "/tmp/partage.sqlite3",
	}
}
