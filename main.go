package main

import (
	"fmt"
	"log"
	"net/http"

	"MorseApp/handlers"
)

func main() {
	// Инициализация маршрутов
	handlers.SetupRoutes()

	// Запуск сервера
	port := ":8080"
	fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
	fmt.Println("📡 Переводчик азбуки Морзе готов к работе!")
	fmt.Println("✨ Откройте браузер и перейдите по адресу выше")

	log.Fatal(http.ListenAndServe(port, nil))
}
