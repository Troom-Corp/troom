// Пакет для работы с jwt tokens
package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// KEY Слово-секрет, нужен для расшифровки токена
var KEY = []byte("secret")

// TOKEN_TIME_ACCESS Время жизни access токена, срок годности
var TOKEN_TIME_ACCESS int64 = 100

// TOKEN_TIME_REFRESH Время жизни refresh токена, срок годности
var TOKEN_TIME_REFRESH int64 = 432000

// CreateAccessToken Метод создания access токена
func CreateAccessToken(userId int) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Создаем payload структуру
		"userId":      userId,                                       // UserId для идентификации пользователя
		"expiredTime": int64(time.Now().Unix()) + TOKEN_TIME_ACCESS, // expiredTime для безопасности
	}).SignedString(KEY)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateRefreshToken Метод создания refrech токена
func CreateRefreshToken(userId int) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Создаем payload структуру
		"userId":      userId,                                        // UserId для идентификации пользователя
		"expiredTime": int64(time.Now().Unix()) + TOKEN_TIME_REFRESH, // expiredTime для безопасности
	}).SignedString(KEY)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetIdentity Расшифровываем токен и получаем из него данные (identity)
func GetIdentity(token string) (int, int64, error) {
	identity, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})
	if err != nil {
		return 0, 0, err
	}

	payload := identity.Claims.(jwt.MapClaims)
	userId := int(payload["userId"].(float64))
	expiredTime := int64(payload["expiredTime"].(float64))

	// Возвращаем payload пользователя в удобных типах данных
	return userId, expiredTime, nil
}
