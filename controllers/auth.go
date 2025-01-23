package controllers

import (
	"encoding/json"
	"forum/config"
	"forum/models"
	"forum/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Password = utils.HashPassword(user.Password)
	config.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var foundUser models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	config.DB.Where("username = ?", user.Username).First(&foundUser)
	if foundUser.ID == 0 || !utils.CheckPassword(foundUser.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Subject:   foundUser.Username,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenString)
}
