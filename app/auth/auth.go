package auth

import (
	"net/http"
	"strings"

	"gorm.io/gorm"

	controller "spender/v1/api/app/controller"
	utils "spender/v1/api/app/utils"

)

var error = utils.CustomError{}

func CheckAuth(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path != "/v1/auth/signup" && path != "/v1/auth/login" {
		
		//required auth
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) < 2 {
			error.ApiError(w, http.StatusForbidden, "Token not provided!")
			return
		}

		token := bearerToken[1]

		_, err := utils.VerifyJwtToken(token)
		if err != nil {
			error.ApiError(w, http.StatusForbidden, err.Error())
			return
		}

		//get data in header and save to user's table
		os := r.Header.Get("os")
		osVersion := r.Header.Get("os_version")
		deviceId := r.Header.Get("device_id")
		deviceModel := r.Header.Get("device_model")
		appId := r.Header.Get("app_id")
		appVersion := r.Header.Get("app_version")
		userId := r.Header.Get("user_id")


		//find user first
		user := controller.GetUserByUUIDOr404(db, userId, w, r)
		if user == nil {
			return
		}
		user.DeviceId = deviceId
		user.DeviceModel = deviceModel
		user.OS = os
		user.OSVersion  = osVersion
		user.AppId = appId
		user.AppVersion = appVersion

		//update user data in the db
		if err := db.Save(&user).Error; err != nil {
		//respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		}
		next.ServeHTTP(w, r)
	})
}