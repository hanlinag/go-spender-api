package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port 	 string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     "localhost",
			Username: "spender_user",
			Password: "root",
			Name:     "spenderdb",
			Port:  	  "5432",
		},
	}
}
