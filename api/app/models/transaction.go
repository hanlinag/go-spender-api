package models

import (


	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Uuid		string		`sql:"type:uuid;primary_key;default:uuid_generate_v4()" json:"uuid"`
	Title   	string      `json:"title"`
	WalletId   	string      `json:"wallet_id"`
	Amount    	string     	`json:"amount"`
	Category 	string   	`json:"category"`
	Type 		string 	    `json:"type"`
	Note 		string		`json:"note"`
	Date 		string		`json:"date"`
	UserId		string		`json:"user_id"`
}


