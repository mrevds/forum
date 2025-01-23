package controllers

import (
	"encoding/json"
	"forum/config"
	"forum/models"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// CreatePost функция для создания нового поста

// CreatePost функция для создания нового поста
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовков
	tokenString := r.Header.Get("Authorization")

	// Удаляем префикс "Bearer " если он присутствует
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:] // Убираем "Bearer "
	} else {
		log.Printf("Invalid Authorization header format\n")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims := &jwt.StandardClaims{}

	// Пытаемся распарсить токен
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Ваш секрет
	})
	if err != nil {
		log.Printf("Unauthorized: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Декодируем тело запроса в структуру поста
	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Printf("Invalid request payload: %v\n", err)
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из токена
	post.UserID = findUserID(claims.Subject)
	log.Printf("User ID obtained: %d\n", post.UserID)

	// Проверяем валидность UserID
	if post.UserID == 0 {
		log.Printf("Invalid user ID\n")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Логируем данные поста перед созданием
	log.Printf("Creating post: Title=%s, Content=%s, UserID=%d\n", post.Title, post.Content, post.UserID)

	// Пытаемся создать пост в базе данных
	if err := config.DB.Create(&post).Error; err != nil {
		log.Printf("Database create error: %s\n", err)
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный пост с ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// findUserID вспомогательная функция для нахождения ID пользователя по его имени
func findUserID(username string) uint {
	var user models.User
	config.DB.Where("username = ?", username).First(&user)
	log.Printf("Found user: %v\n", user) // Логируем пользователя
	return user.ID
}

// DeletePost функция для удаления поста
func DeletePost(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := &jwt.StandardClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	var post models.Post
	config.DB.Delete(&post, vars["id"])
	w.WriteHeader(http.StatusNoContent)
}
