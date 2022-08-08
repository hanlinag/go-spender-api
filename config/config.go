package config

import (
	"time"
)

var ISLOCAL = false

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     "localhost",
			Username: "spender_user",
			Password: "root",
			Name:     "spenderdb",
			Port:     "5432",
		},
	}
}

var JWT_SECRET = "SpEnDerWeBaPiWiThGoLaNg2022"
var TokenExpiredTime = time.Now().Add(60 * time.Minute)
