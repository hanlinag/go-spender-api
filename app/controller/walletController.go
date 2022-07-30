package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	model "spender/v1/app/models"
	models "spender/v1/app/models/Responses"
)

func GetAllWallets(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	wallet := []model.Wallet{}

	db.Where(&model.Wallet{UserId: r.Header.Get("user_id")}).Order("created_at desc").Find(&wallet)

	msg := "Wallet data found"
	if len(wallet) == 0 {
		msg = "No Wallet for this user yet."
	}
	respondJSONWithFormat(w, http.StatusOK, wallet, nil, 200, msg)
}

func CreateWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	wallet := model.Wallet{}

	userId := r.Header.Get("user_id")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wallet); err != nil {

		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")
		return
	}
	defer r.Body.Close()

	wallet.UserId = userId
	wallet.Uuid = uuid.New().String()

	if err := db.Save(&wallet).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving data to the database. Please try again.")
		return
	}

	respondJSONWithFormat(w, http.StatusCreated, wallet, nil, 201, "Data created successfully.")

}

func GetWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	wallet := getWalletOr404(db, uuid, w)
	if wallet == nil {
		return
	}
	respondJSONWithFormat(w, http.StatusOK, wallet, nil, 201, "Data found. ")

}

func UpdateWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	wallet := getWalletOr404(db, uuid, w)
	if wallet == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&wallet); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")
		return
	}
	defer r.Body.Close()

	if err := db.Save(&wallet).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving data. Please try again.")

		return
	}
	respondJSONWithFormat(w, http.StatusOK, wallet, nil, 200, "Data updated successfully.")

}

func DeleteWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	wallet := getWalletOr404(db, uuid, w)
	if wallet == nil {
		return
	}
	if err := db.Delete(&wallet).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error deleting data in the database. Please try again.")

		return
	}

	respondJSONWithFormat(w, http.StatusOK, nil, nil, 204, "Deleted transaciton successfully.")
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getWalletOr404(db *gorm.DB, uuid string, w http.ResponseWriter) *model.Wallet {
	wallet := model.Wallet{}
	if err := db.First(&wallet, model.Wallet{Uuid: uuid}).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "Transaction data not found.")
		return nil
	}
	return &wallet
}
