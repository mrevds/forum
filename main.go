package main

import (
	"fmt" // Импортируем пакет fmt для вывода в консоль
	"forum/config"
	"forum/routes"
	"log" // Импортируем пакет log для логирования ошибок
	"net/http"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRouter()

	fmt.Println("Запуск сервера на порту :8000") // Сообщение о запуске

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err) // Логирование ошибки, если сервер не запустился
	}

	fmt.Println("Сервер успешно остановлен") // Сообщение о штатной остановке сервера (добавлено для полноты картины)
}
