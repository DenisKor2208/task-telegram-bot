package sessionutils

import (
	"fmt"
	"strconv"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/pkg/errors"
)

// ExtractValueFromSession извлекает значение из сессии по ключу и возвращает его как interface{}
func ExtractValueFromSession(session *models.UserSession, key string) (any, error) {
	value, exists := session.Data[key]
	if !exists {
		return nil, errors.New("значение не найдено в сессии")
	}

	return value, nil
}

// ExtractInt64FromSession извлекает значение из сессии и преобразует его в int64
func ExtractInt64FromSession(session *models.UserSession, key string) (int64, error) {
	value, err := ExtractValueFromSession(session, key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		parsedValue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.New("не удалось преобразовать строку в число")
		}
		return parsedValue, nil
	default:
		return 0, errors.New("неподдерживаемый тип данных")
	}
}

// ExtractStringFromSession извлекает значение из сессии и преобразует его в строку
func ExtractStringFromSession(session *models.UserSession, key string) (string, error) {
	value, err := ExtractValueFromSession(session, key)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", errors.New("неподдерживаемый тип данных")
	}
}

// ExtractBoolFromSession извлекает значение из сессии и преобразует его в bool
func ExtractBoolFromSession(session *models.UserSession, key string) (bool, error) {
	value, err := ExtractValueFromSession(session, key)
	if err != nil {
		return false, err
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int, int32, int64, float32, float64:
		// Любое ненулевое значение считается true
		return v != 0, nil
	default:
		return false, errors.New("неподдерживаемый тип данных")
	}
}
