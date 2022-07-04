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

var JWT_SECRET 			= "SpEnDerWeBaPiWiThGoLaNg2022"
var DBURL = "postgres://yzimggivagyaph:8cffacdcd5d764252d3aefe2290707ec24914af7b4cea98f234af94d04d8514d@ec2-52-208-164-5.eu-west-1.compute.amazonaws.com:5432/d32htuv5q4469t"
var TokenExpiredTime  = time.Now().Add(60 * time.Minute)