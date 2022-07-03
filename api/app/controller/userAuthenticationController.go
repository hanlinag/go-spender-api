package controller

import (
	"encoding/json"
	"net/http"

	//"github.com/gorilla/mux"
	"gorm.io/gorm"

	model "spender/v1/api/app/models"
	utils "spender/v1/api/app/utils"
)

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	login := model.Login{}

	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	//find user first
	user := getUserOr404(db, login.Email, w, r)
	if user == nil {
		return
	}

	var token string
	var err error
	//check pw
	if login.Password == user.Password {
		//payload := utils.Payload{}

		//login successs, generate token
		token, err = utils.GenerateJwtToken(utils.Payload{Name: user.Name, Email: user.Email, Id: user.ID,})
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user.Token = token
	} else {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return 
	}
	

	respondJSON(w, http.StatusOK, user)
}



func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, user)
}

func Logout(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//clear token
	respondJSON(w, http.StatusOK, "Successly logged out")
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getUserOr404(db *gorm.DB, email string, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{Email: email}).Error; err != nil {
		respondError(w, http.StatusNotFound, "There is no such user in our system. Please check your credentials and try again.")
		return nil
	}
	return &user
}