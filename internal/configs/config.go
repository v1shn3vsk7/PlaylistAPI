package configs

type Config struct {
	BindAddr string
	LogLevel string
	DbUrl    string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "8888",
		LogLevel: "debug",
		DbUrl:    "postgres://user:password@db:5432/Playlist?sslmode=disable",
	}
}
