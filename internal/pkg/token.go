package pkg

import (
	"github.com/golang-jwt/jwt/v5"
)

// KEY Слово-секрет, нужен для расшифровки токена
var KEY = []byte("secret")

// SignJWT signs a new jwt token with payload
func SignJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": userID,
	})
	accessToken, err := token.SignedString(KEY)
	return accessToken, err
}

// GetIdentity Расшифровываем токен и получаем из него данные (identity)
func GetIdentity(token string) (int, error) {
	identity, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return KEY, nil
	})

	if err != nil {
		return 0, err
	}

	payload := identity.Claims.(jwt.MapClaims)
	userId := int(payload["ID"].(float64))

	return userId, nil

}
