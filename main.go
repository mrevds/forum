package main

import (
	"fmt"
	"forum/config"
	"forum/routes"
	"log"
	"net/http"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRouter()

	fmt.Println("Запуск сервера на порту :8000")

	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}

	fmt.Println("Сервер успешно остановлен")
}
