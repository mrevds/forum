package routes

import (
	"forum/controllers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/posts", controllers.CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", controllers.DeletePost).Methods("DELETE")
	return r
}
