package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
	"time"
)

type AuthControllers struct {
	UserServices store.InterfaceUser
}

func (a AuthControllers) SignIn(c *fiber.Ctx) error {
	var credentials models.SignInCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return fiber.NewError(400, "Bad request")
	}

	user, _ := a.UserServices.IsUserExist(credentials.Login)

	if err := pkg.Decode([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return fiber.NewError(404, "Неверные данные пользователя")
	}

	token, err := pkg.SignJWT(user.UserId)
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании JWT токена")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: "None",
		Expires:  time.Now().Add(time.Minute * 10),
		Secure:   false,
	})

	return fiber.NewError(200, "Вы успешно вошли в аккаунт")
}

func (a AuthControllers) SignUp(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	err := c.BodyParser(&credentials)

	if err != nil {
		return fiber.NewError(400, "Bad request")
	}

	newUserCredentials := models.User{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Login:     credentials.Login,
		Email:     credentials.Email,
		Password:  credentials.Password,
	}

	if !credentials.Validate() {
		return fiber.NewError(400, "Bad request")
	}

	insertedID, err := a.UserServices.InsertOne(newUserCredentials)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(500, "Ошибка при создании пользователя")
	}

	token, err := pkg.SignJWT(insertedID)
	if err != nil {
		return fiber.NewError(500, "Ошибка при создании JWT токена")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: "None",
		Expires:  time.Now().Add(time.Minute * 10),
		Secure:   false,
	})

	return fiber.NewError(201, "Пользователь успешно создан")
}

func (a AuthControllers) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
	})

	return fiber.NewError(200, "Вы успешно вышли из аккаунта")
}

// ValidateCredentials runs on the client before the SignUp method
func (a AuthControllers) ValidateCredentials(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	var isCredentialsValid models.IsCredentials
	c.BodyParser(&credentials)

	isValid, _ := a.UserServices.ValidateCredentials(credentials.Login, credentials.Email)

	for i := range isValid {
		if isValid[i].Login == credentials.Login {
			isCredentialsValid.Login = "Пользователь с таким логином уже существует"
		}
		if isValid[i].Email == credentials.Email {
			isCredentialsValid.Email = "Пользователь с такой почтой уже сущесвует"
		}
	}

	if !pkg.ValidLogin(credentials.Login) {
		isCredentialsValid.Login = "Логин должен соответствовать требованиям"
	}

	if isCredentialsValid.Email != "" || isCredentialsValid.Login != "" {
		invalidMsg, _ := json.Marshal(&isCredentialsValid)
		return fiber.NewError(409, string(invalidMsg))
	}

	return c.JSON(isCredentialsValid)
}

func (u AuthControllers) Profile(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, err := pkg.GetIdentity(jwt)
	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	userProfile, err := u.UserServices.FindOne("userid", ID)

	if err != nil {
		return fiber.NewError(500, "Ошибка при получении профиля")
	}

	return c.JSON(userProfile)
}

func GetAuthControllers(store store.InterfaceStore) *AuthControllers {
	return &AuthControllers{
		UserServices: store.Users(),
	}
}
