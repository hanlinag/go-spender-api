package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"

	model "spender/v1/app/models"
	models "spender/v1/app/models/Responses"
)

func GetAllFeedbacks(db *gorm.DB, w http.ResponseWriter) {
	feedbacks := []model.Feedback{}

	db.Order("created_at desc").Find(&feedbacks)

	msg := "Feedback data found."
	if len(feedbacks) == 0 {
		msg = "No feedback found."
	}
	//msg = ("len %d", len(transatransactions) )
	respondJSONWithFormat(w, http.StatusOK, feedbacks, nil, 200, msg)
}

func CreateFeedback(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	feedback := model.Feedback{}

	userId := r.Header.Get("user_id")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&feedback); err != nil {

		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 400, "Bad request. Please try again.")
		return
	}
	defer r.Body.Close()

	feedback.UserId = userId
	feedback.Uuid = uuid.New().String()

	if err := db.Save(&feedback).Error; err != nil {
		errMsg := &models.ErrorResponse{}
		errMsg.Message = err.Error()

		respondJSONWithFormat(w, http.StatusOK, nil, errMsg, 500, "Error saving data to the database. Please try again.")
		return
	}

	respondJSONWithFormat(w, http.StatusCreated, feedback, nil, 201, "Data created successfully.")

}
