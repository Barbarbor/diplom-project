package i18n

import (
	"bytes"
	"html/template"
	"strings"
)

// Translations хранит переводы для разных языков.
var Translations = map[string]map[string]map[string]map[string]string{
	"en": {
		"survey": {
			"handler": {
				"notFound":     "Survey not found",
				"accessDenied": "Access denied",
				"success":      "Survey created successfully",
				"invalidData":  "Invalid survey data",
			},
			"service": {
				"defaultTitle":  "Survey from {{.Date}}",
				"hashError":     "Failed to generate survey hash", // Internal, untranslated
				"creationError": "Failed to create survey",        // Internal, untranslated
			},
			"repository": {
				"notFound":      "Survey not found in database", // Internal, untranslated
				"creationError": "Failed to create survey",      // Internal, untranslated
			},
		},
		"question": {
			"handler": {
				"notFound":    "Question not found",
				"success":     "Question created successfully",
				"invalidType": "Invalid question type",
				"invalidData": "Invalid question data",
			},
			"service": {
				"defaultSingle":    "Choose one option",
				"defaultMulti":     "Choose multiple options",
				"defaultConsent":   "Please provide your consent",
				"defaultEmail":     "Enter your email",
				"defaultRating":    "Please rate on the scale",
				"defaultDate":      "Select a date",
				"defaultShortText": "Short answer",
				"defaultLongText":  "Long answer",
				"defaultNumber":    "Enter a number",
				"invalidType":      "Invalid question type", // Internal, untranslated
			},
		},
		"option": {
			"handler": {
				"notFound":               "Option not found",
				"invalidOptionID":        "Invalid option ID",
				"invalidQuestionContext": "Invalid question context",
			},
		},
	},
	"ru": {
		"survey": {
			"handler": {
				"notFound":     "Опрос не найден",
				"accessDenied": "Доступ запрещён",
				"success":      "Опрос успешно создан",
				"invalidData":  "Неверные данные опроса",
			},
			"service": {
				"defaultTitle":  "Опрос от {{.Date}}",
				"hashError":     "Не удалось сгенерировать хэш опроса", // Internal, untranslated
				"creationError": "Не удалось создать опрос",            // Internal, untranslated
			},
			"repository": {
				"notFound":      "Опрос не найден в базе данных", // Internal, untranslated
				"creationError": "Не удалось создать опрос",      // Internal, untranslated
			},
		},
		"question": {
			"handler": {
				"notFound":    "Вопрос не найден",
				"success":     "Вопрос успешно создан",
				"invalidType": "Неверный тип вопроса",
				"invalidData": "Неверные данные вопроса",
			},
			"service": {
				"defaultSingle":    "Выберите один вариант",
				"defaultMulti":     "Выберите несколько вариантов",
				"defaultConsent":   "Пожалуйста, подтвердите своё согласие",
				"defaultEmail":     "Введите ваш email",
				"defaultRating":    "Пожалуйста, оцените по шкале",
				"defaultDate":      "Выберите дату",
				"defaultShortText": "Краткий ответ",
				"defaultLongText":  "Развернутый ответ",
				"defaultNumber":    "Введите число",

				"invalidType": "Неверный тип вопроса", // Internal, untranslated
			},
		},
		"option": {
			"handler": {
				"notFound":               "Опция не найдена",
				"invalidOptionID":        "Неверный идентификатор опции",
				"invalidQuestionContext": "Неправильный контекст вопроса",
			},
		},
	},
}

// LanguageProvider defines an interface for retrieving the current language
type LanguageProvider interface {
	GetLang() string
}

// RequestLanguageProvider implements LanguageProvider for request-scoped language
type RequestLanguageProvider struct {
	Lang string // Exported field (uppercase L)
}

func (p *RequestLanguageProvider) GetLang() string {
	if p.Lang == "" {
		return "ru" // Default language
	}
	return p.Lang
}

// Global language provider instance
var langProvider LanguageProvider

// SetLanguageProvider sets the current language provider
func SetLanguageProvider(provider LanguageProvider) {
	langProvider = provider
}

// T возвращает перевод для заданного ключа
func T(key string) string {
	if langProvider == nil {
		return key // Fallback if no provider is set
	}
	lang := langProvider.GetLang()
	parts := strings.SplitN(key, ".", 3)
	if len(parts) != 3 {
		return key
	}
	entity, layer, keyName := parts[0], parts[1], parts[2]

	if entityMap, ok := Translations[lang]; ok {
		if layerMap, ok := entityMap[entity]; ok {
			if keyMap, ok := layerMap[layer]; ok {
				if val, ok := keyMap[keyName]; ok {
					return val
				}
			}
		}
	}
	return key
}

// TWithData возвращает перевод с подстановкой данных
func TWithData(key string, data interface{}) (string, error) {
	if langProvider == nil {
		return key, nil // Fallback if no provider is set
	}
	lang := langProvider.GetLang()
	parts := strings.SplitN(key, ".", 3)
	if len(parts) != 3 {
		return key, nil
	}
	entity, layer, keyName := parts[0], parts[1], parts[2]

	if entityMap, ok := Translations[lang]; ok {
		if layerMap, ok := entityMap[entity]; ok {
			if keyMap, ok := layerMap[layer]; ok {
				if val, ok := keyMap[keyName]; ok {
					tmpl, err := template.New("translation").Parse(val)
					if err != nil {
						return val, err
					}
					var buf bytes.Buffer
					err = tmpl.Execute(&buf, data)
					if err != nil {
						return val, err
					}
					return buf.String(), nil
				}
			}
		}
	}
	return key, nil
}
