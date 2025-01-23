package controllers

import (
	"encoding/json"
	"forum/config"
	"forum/models"
	"log"
	"net/http"
	"strconv" // Import strconv for string to int conversion

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// CreatePost функция для создания нового поста
func CreatePost(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Println("Missing Authorization header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := validateToken(tokenString) // Use validateToken function
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest) // Simplified error message
		return
	}

	post.UserID = findUserID(claims.Subject)
	if post.UserID == 0 {
		log.Println("Invalid user ID")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Creating post: Title=%s, Content=%s, UserID=%d\n", post.Title, post.Content, post.UserID)

	if err := config.DB.Create(&post).Error; err != nil {
		log.Printf("Database create error: %s\n", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError) // Simplified error message
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// findUserID вспомогательная функция для нахождения ID пользователя по его имени
func findUserID(username string) uint {
	var user models.User
	config.DB.Where("username = ?", username).First(&user)
	log.Printf("Found user: %v\n", user)
	return user.ID
}

// DeletePost функция для удаления поста
func DeletePost(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Println("Missing Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := validateToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"]) // Convert string ID to integer
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	result := config.DB.Where("id = ? AND user_id = ?", postID, findUserID(claims.Subject)).Delete(&post) // Delete with condition
	if result.Error != nil {
		log.Printf("Database delete error: %v", result.Error)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("Post with ID %d not found or you are not the owner", postID)
		http.Error(w, "Post not found or you are not the owner", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Post with ID %d successfully deleted", postID)
}

func validateToken(tokenString string) (*jwt.StandardClaims, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		return nil, jwt.ErrSignatureInvalid
	}

	claims := &jwt.StandardClaims{}
	secret := config.Cfg.JWT.Secret // Retrieve secret from configuration
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
