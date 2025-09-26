package handlers

import (
	"net/http"
)

// SetupRoutes настраивает маршруты приложения
func SetupRoutes() {
	// Статические файлы
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Маршруты
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/translate", TranslateHandler)
}
