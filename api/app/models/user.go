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
	LastLogin		time.Time		`json:"last_login"`
	IsLogin 		bool  			 `json:"is_login"`
	DeviceId		string			`json:"device_id"`
	DeviceModel		string			`json:"device_model"`
	OS				string			`json:"os"`
	OSVersion		string			`json:"os_version"`
	AppId			string 			`json:"app_id"`
	AppVersion		string			`json:"app_version"`

}

func (u *User) SetLogin(isLogin bool) {
	u.IsLogin = isLogin
}

func (u *User) SetActive(status bool) {
	u.IsActive = status
}

func (u *User) SetVerified(status bool) {
	u.IsVerified = status
}
