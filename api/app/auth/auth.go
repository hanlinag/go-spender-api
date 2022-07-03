package auth

import (
	"net/http"
	"strings"

	"spender/v1/api/app/utils"
)

var error = utils.CustomError{}

func CheckAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r)
	})
}