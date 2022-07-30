package models

type GeneralResponse struct {
	StatusCode 		int `json:"status_code"`
	Desc 	 		string `json:"description"`
	Timestamp 		string `json:"timestamp"`
	Error 			interface{} `json:"error"`
	Data 			interface{} `json:"data"`
}