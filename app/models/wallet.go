package models

import (
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	Uuid   string `sql:"type:uuid;primary_key;default:uuid_generate_v4()" json:"uuid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Date   string `json:"date"`
	UserId string `json:"user_id"`
}
