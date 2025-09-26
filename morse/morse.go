package morse

import (
	"strings"
	"unicode"
)

// Карта символов в азбуку Морзе
var morseCode = map[rune]string{
	// Латинские буквы
	'A': ".-", 'B': "-...", 'C': "-.-.", 'D': "-..", 'E': ".", 'F': "..-.",
	'G': "--.", 'H': "....", 'I': "..", 'J': ".---", 'K': "-.-", 'L': ".-..",
	'M': "--", 'N': "-.", 'O': "---", 'P': ".--.", 'Q': "--.-", 'R': ".-.",
	'S': "...", 'T': "-", 'U': "..-", 'V': "...-", 'W': ".--", 'X': "-..-",
	'Y': "-.--", 'Z': "--..",

	// Русские буквы
	'А': ".-", 'Б': "-...", 'В': ".--", 'Г': "--.", 'Д': "-..", 'Е': ".",
	'Ё': ".", 'Ж': "...-", 'З': "--..", 'И': "..", 'Й': ".---", 'К': "-.-",
	'Л': ".-..", 'М': "--", 'Н': "-.", 'О': "---", 'П': ".--.", 'Р': ".-.",
	'С': "...", 'Т': "-", 'У': "..-", 'Ф': "..-.", 'Х': "....", 'Ц': "-.-.",
	'Ч': "---.", 'Ш': "----", 'Щ': "--.-", 'Ъ': "--.--", 'Ы': "-.--",
	'Ь': "-..-", 'Э': "..-..", 'Ю': "..--", 'Я': ".-.-",

	// Цифры
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-",
	'5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",

	// Знаки препинания
	'.': ".-.-.-", ',': "--..--", '?': "..--..", '\'': ".----.", '!': "-.-.--",
	'/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...",
	';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-",
	'"': ".-..-.", '$': "...-..-", '@': ".--.-.", ' ': "/",
}

// Обратные карты для декодирования с приоритетом русского языка
var (
	textFromMorseRU map[string]rune // Для русского языка
	textFromMorseEN map[string]rune // Для английского языка
)

func init() {
	// Инициализируем карты
	textFromMorseRU = make(map[string]rune)
	textFromMorseEN = make(map[string]rune)

	// Заполняем карты, отдавая приоритет русским буквам при конфликтах
	for char, morse := range morseCode {
		// Определяем, русская это буква или английская
		if (char >= 'А' && char <= 'Я') || char == 'Ё' {
			textFromMorseRU[morse] = char
		} else if char >= 'A' && char <= 'Z' {
			textFromMorseEN[morse] = char
		} else {
			// Для цифр и знаков препинания добавляем в обе карты
			textFromMorseRU[morse] = char
			textFromMorseEN[morse] = char
		}
	}
}

// TextToMorse преобразует текст в азбуку Морзе
func TextToMorse(text string) string {
	text = strings.ToUpper(text)
	var result []string

	for _, char := range text {
		if morse, exists := morseCode[char]; exists {
			result = append(result, morse)
		} else if char == '\n' {
			result = append(result, "\n")
		} else if unicode.IsSpace(char) {
			result = append(result, "/")
		} else {
			result = append(result, "..--..")
		}
	}

	return strings.Join(result, " ")
}

// MorseToText преобразует азбуку Морзе в текст с определением языка
func MorseToText(morse string) string {
	morse = strings.ReplaceAll(morse, "  ", " ")
	morse = strings.TrimSpace(morse)

	if morse == "" {
		return ""
	}

	// Определяем вероятный язык на основе всего текста
	language := detectLanguageFromMorse(morse)

	var result []string
	lines := strings.Split(morse, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			result = append(result, "")
			continue
		}

		words := strings.Split(line, "/")
		var decodedWords []string

		for _, word := range words {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}

			letters := strings.Split(word, " ")
			var decodedLetters []string

			for _, letter := range letters {
				letter = strings.TrimSpace(letter)
				if letter == "" {
					continue
				}

				var decodedChar rune
				var found bool

				// Используем карту в зависимости от определенного языка
				if language == "русский" {
					decodedChar, found = textFromMorseRU[letter]
					if !found {
						decodedChar, found = textFromMorseEN[letter]
					}
				} else {
					decodedChar, found = textFromMorseEN[letter]
					if !found {
						decodedChar, found = textFromMorseRU[letter]
					}
				}

				if found {
					decodedLetters = append(decodedLetters, string(decodedChar))
				} else {
					decodedLetters = append(decodedLetters, "?")
				}
			}

			if len(decodedLetters) > 0 {
				decodedWords = append(decodedWords, strings.Join(decodedLetters, ""))
			}
		}

		result = append(result, strings.Join(decodedWords, " "))
	}

	return strings.Join(result, "\n")
}

// MorseToTextWithLanguage преобразует азбуку Морзе в текст с указанием языка
func MorseToTextWithLanguage(morse string, language string) string {
	morse = strings.ReplaceAll(morse, "  ", " ")
	morse = strings.TrimSpace(morse)

	if morse == "" {
		return ""
	}

	var result []string
	lines := strings.Split(morse, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			result = append(result, "")
			continue
		}

		words := strings.Split(line, "/")
		var decodedWords []string

		for _, word := range words {
			word = strings.TrimSpace(word)
			if word == "" {
				continue
			}

			letters := strings.Split(word, " ")
			var decodedLetters []string

			for _, letter := range letters {
				letter = strings.TrimSpace(letter)
				if letter == "" {
					continue
				}

				var decodedChar rune
				var found bool

				// Используем выбранный язык
				if language == "english" {
					decodedChar, found = textFromMorseEN[letter]
					if !found {
						decodedChar, found = textFromMorseRU[letter]
					}
				} else {
					// По умолчанию русский
					decodedChar, found = textFromMorseRU[letter]
					if !found {
						decodedChar, found = textFromMorseEN[letter]
					}
				}

				if found {
					decodedLetters = append(decodedLetters, string(decodedChar))
				} else {
					decodedLetters = append(decodedLetters, "?")
				}
			}

			if len(decodedLetters) > 0 {
				decodedWords = append(decodedWords, strings.Join(decodedLetters, ""))
			}
		}

		result = append(result, strings.Join(decodedWords, " "))
	}

	return strings.Join(result, "\n")
}

// detectLanguageFromMorse пытается определить язык по коду Морзе
func detectLanguageFromMorse(morse string) string {
	// Считаем количество русских и английских букв
	var russianCount, englishCount int

	words := strings.Split(morse, "/")
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		letters := strings.Split(word, " ")
		for _, letter := range letters {
			letter = strings.TrimSpace(letter)
			if letter == "" {
				continue
			}

			// Проверяем, существует ли буква только в русской карте
			_, ruExists := textFromMorseRU[letter]
			_, enExists := textFromMorseEN[letter]

			// Если буква существует только в русской карте - увеличиваем счетчик русского
			if ruExists && !enExists {
				russianCount++
			} else if enExists && !ruExists {
				// Если буква существует только в английской карте - увеличиваем счетчик английского
				englishCount++
			}
			// Если буква существует в обеих картах (конфликт), не увеличиваем счетчики
		}
	}

	if russianCount > englishCount {
		return "русский"
	} else if englishCount > russianCount {
		return "английский"
	}

	// По умолчанию считаем русским (так как чаще используется русский текст)
	return "русский"
}

// IsValidMorse проверяет, является ли строка корректным кодом Морзе
func IsValidMorse(text string) bool {
	if strings.TrimSpace(text) == "" {
		return false
	}

	for _, char := range text {
		if char != '.' && char != '-' && char != ' ' && char != '/' && char != '\n' {
			return false
		}
	}

	return true
}

// DetectLanguage определяет язык текста (для текстового ввода)
func DetectLanguage(text string) string {
	var russianCount, englishCount int

	for _, char := range text {
		if unicode.IsLetter(char) {
			if (char >= 'А' && char <= 'Я') || (char >= 'а' && char <= 'я') || char == 'Ё' || char == 'ё' {
				russianCount++
			} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
				englishCount++
			}
		}
	}

	if russianCount > englishCount {
		return "русский"
	} else if englishCount > russianCount {
		return "английский"
	}
	return "неопределен"
}
