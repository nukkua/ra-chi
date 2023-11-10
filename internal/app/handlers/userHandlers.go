package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nukkua/ra-chi/internal/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		db.Find(&users)
		json.NewEncoder(w).Encode(users)
	}
}

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hashedPassword, errorHashing := bcrypt.GenerateFromPassword(user.PasswordHash, bcrypt.MinCost)
		if errorHashing != nil {
			http.Error(w, errorHashing.Error(), http.StatusInternalServerError)
		}

		user.PasswordHash = hashedPassword
		result := db.Create(&user)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
