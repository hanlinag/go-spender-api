package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	model "spender/v1/app/models"
	models "spender/v1/app/models/Responses"
	utils "spender/v1/app/utils"
)

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	login := model.Login{}

	response := &models.GeneralResponse{}
	response.Timestamp = time.Now().Local().String()

	os := r.Header.Get("os")
	osVersion := r.Header.Get("os_version")
	deviceId := r.Header.Get("device_id")
	deviceModel := r.Header.Get("device_model")
	appId := r.Header.Get("app_id")
	appVersion := r.Header.Get("app_version")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad Request. Please check and try again.")
		return
	}
	defer r.Body.Close()

	//find user first
	user := GetUserOr404(db, login.Email, w, r)
	if user == nil {
		return
	}

	var token string
	var err error
	//check pw
	if login.Password == user.Password {
		//payload := utils.Payload{}

		//check if the current account is logged in yet?
		if user.IsLogin && deviceId != user.DeviceId {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = "This account has been logged in with another device. Please log out from other device first and try again."

			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "You have loggined this account in other account. Please log out from the other device first and try again.")

			return
		}

		//is verified or not
		if !user.IsVerified {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = "This account is not verified yet. Please check your email and verify."

			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "This account is verified yet. Please check your email and verify.")

			return

		}

		//is active or not
		if !user.IsActive {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = "This account is not active anymore. Please contact the support to reactivate your account."

			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "This account is not active anymore. Please contact the support to reactivate your account.")

			return

		}

		//login successs, generate token
		token, err = utils.GenerateJwtToken(utils.Payload{Name: user.Name, Email: user.Email, Id: user.ID})
		if err != nil {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = err.Error()
			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error generating token. Please try again.")

			return
		}
		user.Token = token

		//set last login time and status
		user.SetLogin(true)
		user.LastLogin = time.Now()
		user.DeviceId = deviceId
		user.DeviceModel = deviceModel
		user.OS = os
		user.OSVersion = osVersion
		user.AppId = appId
		user.AppVersion = appVersion

		fmt.Sprintf("Variable string %s content", user.OSVersion)

		//update user data in the db
		if err := db.Save(&user).Error; err != nil {
			errMsg := &models.ErrorResponse{}
			errMsg.Message = err.Error()
			respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving useer data. Please try again.")
			return
		}

	} else {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = "Username or Password incorrect."
		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Unauthorized.")
		return
	}

	respondJSONWithFormat(w, http.StatusOK, user, nil, 200, "Login successfully. ")

}

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	response := &models.GeneralResponse{}
	response.Timestamp = time.Now().Local().String()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad User Input. Please check and try again.")

		return
	}
	defer r.Body.Close()

	user.IsActive = false
	user.IsLogin = false
	user.IsVerified = false

	user.Uuid = uuid.New().String()

	if err := db.Save(&user).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Internal Server Error. Please try again.")

		return
	}

	respondJSONWithFormat(w, http.StatusCreated, user, nil, 201, "Successfully created account.")
}

func Logout(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	deviceId := r.Header.Get("device_id")
	userUUID := r.Header.Get("user_id")

	//find user first
	user := GetUserByUUIDOr404(db, userUUID, w, r)
	if user == nil {
		return
	}

	user.DeviceId = deviceId
	user.IsLogin = false

	UpdateUserDataAfterLogout(db, user, w, r)
}

func UpdateUserDataAfterLogout(db *gorm.DB, user *model.User, w http.ResponseWriter, r *http.Request) {
	response := &models.GeneralResponse{}
	response.Timestamp = time.Now().Local().String()

	if err := db.Save(&user).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "500 Internal Server Error. Please try again later.")

		return
	}

	//success
	respondJSONWithFormat(w, http.StatusOK, nil, nil, 200, "Successfully logged out")

}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]

	user := GetUserByUUIDOr404(db, uuid, w, r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")

		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error updating data. Please try agin.")

		return
	}

	respondJSONWithFormat(w, http.StatusOK, user, nil, 200, "Data updated successfully.")
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func GetUserOr404(db *gorm.DB, email string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Email: email}).Error; err != nil {

		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "There is no such user in our system. Please check your credentials and try again.")
		return nil
	}
	return &user
}

func GetUserByUUIDOr404(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Uuid: uuid}).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "There is no such user in our system. Please check your credentials and try again.")

		return nil
	}
	return &user
}
