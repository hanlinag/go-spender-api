package controller

import (
	"encoding/json"
	"net/http"
	models "spender/v1/app/models/Responses"
	"time"
)

// respondJSON makes the response with payload as json format
func commonRespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
// func respondError(w http.ResponseWriter, code int, message string) {
// 	respondJSON(w, code, map[string]string{"error": message})
// }

// respondError makes the error response with payload as json format
func respondJSONWithFormat(w http.ResponseWriter, code int, data interface{}, error interface{}, customCode int, desc string) {

	response := &models.GeneralResponse{}
	response.Timestamp = time.Now().Local().String()
	response.Desc      = desc 
	response.StatusCode= customCode

	if data != nil {
		response.Data = data
	} 
	
	if error != nil {
		response.Error = error
	}

	commonRespondJSON(w, code, response)
}
