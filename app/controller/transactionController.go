package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	model "spender/v1/app/models"
	models "spender/v1/app/models/Responses"
)

func GetAllTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transactions := []model.Transaction{}

	//get queries
	walletID := r.URL.Query().Get("wallet_id")
	category := r.URL.Query().Get("category")
	transactionType := r.URL.Query().Get("type")
	limitt := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")

	//2021-01-01 00:00:00
	//Format("2006-01-02 15:04:05")

	limit := 1 //DEFAULT LIMIT 10
	if limitt != "" {
		x, err := strconv.ParseInt(limitt, 10, 32)

		if err != nil {
			//	limit = int(x)
		}
		limit = int(x)
	}
	if cursor != "" {

		db.Where("date < ?", cursor).Where(&model.Transaction{UserId: r.Header.Get("user_id"), Type: transactionType, WalletId: walletID, Category: category}).Order("updated_at desc").Limit(limit).Find(&transactions)
	} else {
		db.Where(&model.Transaction{UserId: r.Header.Get("user_id"), Type: transactionType, WalletId: walletID, Category: category}).Order("updated_at desc").Limit(limit).Find(&transactions)
	}

	//db.Where("user_id = ? AND type = ? AND wallet_id = ? AND category = ?", r.Header.Get("user_id"), transactionType, walletID, category).Limit(limit).Find(&transactions)
	//db.Where("user_id = ?", r.Header.Get("user_id")).Find(&transactions)
	//fmt.Println(transactions)
	//to do pagiation
	msg := "Transactions data found"
	if len(transactions) == 0 {
		msg = "No transaciton for this user yet."
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
	transaction := getTransactionOr404(db, uuid, w)
	if transaction == nil {
		return
	}
	respondJSONWithFormat(w, http.StatusOK, transaction, nil, 201, "Data found. ")

}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	transaction := getTransactionOr404(db, uuid, w)
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
	transaction := getTransactionOr404(db, uuid, w)
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
func getTransactionOr404(db *gorm.DB, uuid string, w http.ResponseWriter) *model.Transaction {
	transaction := model.Transaction{}
	if err := db.First(&transaction, model.Transaction{Uuid: uuid}).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 404, "Transaction data not found.")
		return nil
	}
	return &transaction
}
