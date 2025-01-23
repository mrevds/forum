package main

import (
	"net/http"
	"forum/config"
	"forum/routes"
)

func main() {
	config.ConnectDatabase()        // Подключение к базе данных
	r := routes.SetupRouter()       // Настройка маршрутов
	http.ListenAndServe(":8000", r) // Запуск сервера
}
