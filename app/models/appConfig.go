package models
import (
	"gorm.io/gorm"
)

type AppConfig struct {
	gorm.Model
	DataId					string 	`gorm:"unique" json:"data_id"`
	MinVersioniOS 			string `json:"ios_min_version"`
	LatestVersioniOS 		string `json:"ios_latest_version"`
	MinVersionAndroid 		string `json:"android_min_version"`
	LatestVersionAndroid 	string `json:"android_latest_version"`
}