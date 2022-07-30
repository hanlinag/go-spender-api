package controller

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	model "spender/v1/app/models"
	models "spender/v1/app/models/Responses"
	"time"
)

func GetAppConfig(db *gorm.DB, w http.ResponseWriter) {
	appConfig := model.AppConfig{}

	if err := db.First(&appConfig).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "App config data not found")
		return
	}
	respondJSONWithFormat(w, http.StatusOK, appConfig, nil, 200, "App config data found.")

}

func UpdateAppConfig(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	newAppConfig := &model.AppConfig{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newAppConfig); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad Request. Please check and try again.")
		return
	}
	defer r.Body.Close()

	appConfig := model.AppConfig{}
	if err := db.First(&appConfig, model.AppConfig{DataId: newAppConfig.DataId}).Error; err != nil {
		//add new data 
		if err := db.Save(&newAppConfig).Error; err != nil {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = err.Error()
	
			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error Saving App Config Data")
			return
		}
	} else {
		//update ddata
		appConfig.MinVersioniOS = newAppConfig.MinVersioniOS
		appConfig.LatestVersioniOS = newAppConfig.LatestVersioniOS
		appConfig.MinVersionAndroid = newAppConfig.MinVersionAndroid
		appConfig.LatestVersionAndroid = newAppConfig.LatestVersionAndroid

		if err := db.Save(&appConfig).Error; err != nil {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = err.Error()
	
			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error Updating App Config Data")
			return
		}
	}


	respondJSONWithFormat(w, http.StatusOK, newAppConfig, nil, 200, "Data updated successfully.")
}



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
	response.Timestamp = time.Now().UTC().String()
	response.Desc = desc
	response.StatusCode = customCode

	if data != nil {
		response.Data = data
	}

	if error != nil {
		response.Error = error
	}

	commonRespondJSON(w, code, response)
}
