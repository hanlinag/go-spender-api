package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	model "spender/v1/api/app/models"
)

func GetAllTransactions(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transactions := []model.Transaction{}
	db.Where("user_id = ?", r.Header.Get("user_id")).Find(&transactions)
	fmt.Println(transactions)
	respondJSON(w, http.StatusOK, transactions)
}

func CreateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	transaction := model.Transaction{}

	userId := r.Header.Get("user_id")

	
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaction); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	transaction.UserId = userId
	transaction.Uuid = uuid.New().String()

	if err := db.Save(&transaction).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, transaction)
}

func GetTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["user_id"]
	employee := getTransactionOr404(db, name, w, r)
	if employee == nil {
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func UpdateTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, employee)
}

func DeleteTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	employee := getEmployeeOr404(db, name, w, r)
	if employee == nil {
		return
	}
	if err := db.Delete(&employee).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// getEmployeeOr404 gets a employee instance if exists, or respond the 404 error otherwise
func getTransactionOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Transaction {
	transaction := model.Transaction{}
	if err := db.First(&transaction, model.Employee{Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &transaction
}
