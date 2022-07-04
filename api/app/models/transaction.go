package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Uuid		string		`sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title   	string      `json:"title"`
	WalletId   	string      `json:"wallet_id"`
	Amount    	int     	`json:"amount"`
	Category 	string   	`json:"category"`
	Type 		string 	    `json:"type"`
	Note 		string		`json:"note"`
	Date 		time.Time	`json:"date"`
	UserId		string	
}


