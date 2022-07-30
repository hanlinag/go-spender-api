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

func GetAllTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transactions := []model.Transaction{}
	db.Where("user_id = ?", r.Header.Get("user_id")).Find(&transactions)
	//fmt.Println(transactions)
	//to do pagiation
	msg := "Transactions data found"
	if len(transactions) == 0 {
		msg = "No transaciton record for this user yet."
	}
	//msg = ("len %d", len(transatransactions) )
	respondJSONWithFormat(w, http.StatusOK, transactions, nil, 200, msg)
}

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}

	userId := r.Header.Get("user_id")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaction); err != nil {
		//respondError(w, http.StatusBadRequest, err.Error())
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")
		return
	}
	defer r.Body.Close()

	transaction.UserId = userId
	transaction.Uuid = uuid.New().String()

	if err := db.Save(&transaction).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving data to the database. Please try again.")
		return
	}

	respondJSONWithFormat(w, http.StatusCreated, transaction, nil, 201, "Data created successfully.")

}

func GetTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	transaction := getTransactionOr404(db, uuid, w, r)
	if transaction == nil {
		return
	}
	respondJSONWithFormat(w, http.StatusOK, transaction, nil, 201, "Data found. ")

}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	transaction := getTransactionOr404(db, uuid, w, r)
	if transaction == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaction); err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")
		return
	}
	defer r.Body.Close()

	if err := db.Save(&transaction).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving data. Please try again.")

		return
	}
	respondJSONWithFormat(w, http.StatusOK, transaction, nil, 200, "Data updated successfully.")

}

func DeleteTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	transaction := getTransactionOr404(db, uuid, w, r)
	if transaction == nil {
		return
	}
	if err := db.Delete(&transaction).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error deleting data in the database. Please try again.")

		return
	}

	respondJSONWithFormat(w, http.StatusOK, nil, nil, 204, "Deleted transaciton successfully.")
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getTransactionOr404(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *model.Transaction {
	transaction := model.Transaction{}
	if err := db.First(&transaction, model.Transaction{Uuid: uuid}).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "Transaction data not found.")
		return nil
	}
	return &transaction
}
