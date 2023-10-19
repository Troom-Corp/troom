// Пакет для работы с jwt tokens
package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// KEY Слово-секрет, нужен для расшифровки токена
var KEY = []byte("secret")

// TOKEN_TIME Время жизни токена, срок годности
var TOKEN_TIME int64 = 100

// SignJWT Метод создания токена
func SignJWT(userId int) string {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Создаем payload структуру
		"userId":      userId,                                // UserId для идентификации пользователя
		"expiredTime": int64(time.Now().Unix()) + TOKEN_TIME, // expiredTime для безопасности
	}).SignedString(KEY)
	if err != nil {
		return ""
	}
	return token
}

// GetIdentity Расшифровываем токен и получаем из него данные (identity)
func GetIdentity(token string) (int, int64) {
	identity, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})
	if err != nil {
		panic(err)
	}

	payload := identity.Claims.(jwt.MapClaims)
	userId := int(payload["userId"].(float64))
	expiredTime := int64(payload["expiredTime"].(float64))

	// Возвращаем payload пользователя в удобных типах данных
	return userId, expiredTime
}
