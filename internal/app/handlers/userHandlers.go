package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nukkua/ra-chi/internal/app/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtKey = []byte("secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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
		hashedPassword, errorHashing := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if errorHashing != nil {
			http.Error(w, errorHashing.Error(), http.StatusInternalServerError)
			return
		}

		user.PasswordHash = hashedPassword
		user.Password = ""
		result := db.Create(&user)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "user created"})

		// Esto es para devolver un json siempre usado en handlers
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(user)
	}
}
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&creds)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		result := db.Where("username = ?", creds.Username).First(&user)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusUnauthorized)
			return
		}


		error := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(creds.Password));
		if  error != nil {
			http.Error(w, error.Error(), http.StatusUnauthorized)
			return
		}
		expirationTime := time.Now().Add(30 * time.Minute)

		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(JwtKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}
