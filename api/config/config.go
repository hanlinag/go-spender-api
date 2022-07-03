package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host     string
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     "tcp(localhost:3306)",
			Dialect:  "mysql",
			Username: "",
			Password: "",
			Name:     "",
			Charset:  "utf8",
		},
	}
}
