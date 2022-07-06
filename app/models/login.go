package models

type Login struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
	LoginType 	string `json:"login_type"`
}