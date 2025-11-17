package config

type Config struct {
	DBPath string
}

func Load() *Config {
	return &Config{
		DBPath: "./ctf.db", // located in cmd/server/
	}
}
