package models

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	Uuid    string `sql:"type:uuid;primary_key;default:uuid_generate_v4()" json:"uuid"`
	Name    string `json:"name"`
	Rating  string `json:"rating"`
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}
