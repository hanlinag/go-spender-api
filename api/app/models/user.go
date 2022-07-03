package models

import (
	"time"

	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	Uuid	 		string 			`sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name     		string			`json:"name"`
	Email    		string 			`gorm:"unique" json:"email"`
	Nickname 		string 			`json:"nickname"`
	Password 		string			`json:"password"`
	DoB		 		time.Time		`json:"dob"`
	LoginType 		string			`json:"login_type"`
	Occupation 		string			`json:"occupation"`
	IsActive 		bool			`json:"is_active"`
	IsVerified		bool			`json:"is_verified"`
	Token 			string 			`json:"token"`
}

func (u *User) setActive(status bool) {
	u.IsActive = status
}

func (u *User) setVerified(status bool) {
	u.IsVerified = status
}
