package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"

	"MorseApp/morse"
)

// TranslationRequest структура для запроса перевода
type TranslationRequest struct {
	Text     string `json:"text"`
	Mode     string `json:"mode"`
	Language string `json:"language"` // Язык для декодирования
}

// TranslationResponse структура для ответа перевода
type TranslationResponse struct {
	Success  bool   `json:"success"`
	Result   string `json:"result"`
	Language string `json:"language,omitempty"`
	Error    string `json:"error,omitempty"`
}

// HomeHandler обработчик главной страницы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	HTMLTemplate, _ := os.ReadFile("templates/template.html")
	tmpl := template.Must(template.New("index").Parse(string(HTMLTemplate)))
	tmpl.Execute(w, nil)
}

// TranslateHandler API обработчик для перевода
func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		json.NewEncoder(w).Encode(TranslationResponse{
			Success: false,
			Error:   "Метод не поддерживается",
		})
		return
	}

	var req TranslationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(TranslationResponse{
			Success: false,
			Error:   "Ошибка чтения запроса",
		})
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		json.NewEncoder(w).Encode(TranslationResponse{
			Success: false,
			Error:   "Введите текст для перевода",
		})
		return
	}

	var result string
	var language string

	switch req.Mode {
	case "text-to-morse":
		language = morse.DetectLanguage(req.Text)
		result = morse.TextToMorse(req.Text)
	case "morse-to-text":
		// Используем выбранный язык для декодирования
		result = morse.MorseToTextWithLanguage(req.Text, req.Language)
	case "auto":
		if morse.IsValidMorse(req.Text) {
			// Для авторежима используем русский по умолчанию
			result = morse.MorseToTextWithLanguage(req.Text, "russian")
		} else {
			language = morse.DetectLanguage(req.Text)
			result = morse.TextToMorse(req.Text)
		}
	default:
		json.NewEncoder(w).Encode(TranslationResponse{
			Success: false,
			Error:   "Неверный режим перевода",
		})
		return
	}

	json.NewEncoder(w).Encode(TranslationResponse{
		Success:  true,
		Result:   result,
		Language: language,
	})
}
