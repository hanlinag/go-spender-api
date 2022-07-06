package config

import (
	"time"
)

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
var DBURL = "postgres://ibwwxigahwxcgj:61f5148dac4c00f33d833de69771c6580a7d9941cb53781c08ca8eb530248df3@ec2-99-81-137-11.eu-west-1.compute.amazonaws.com:5432/d1llbgccsbfj9h"
var TokenExpiredTime = time.Now().Add(60 * time.Minute)
