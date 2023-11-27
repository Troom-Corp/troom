package controllers

import (
	"encoding/json"
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
)

type AuthControllers struct {
	UserServices store.InterfaceUser
}

func (a *AuthControllers) UserSignIn(c *fiber.Ctx) error {
	var credentials models.SignInCredentials
	c.BodyParser(&credentials)

	user, _ := a.UserServices.UserExists(credentials.Login)

	if err := pkg.Decode([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return fiber.NewError(404, "Неверные данные пользователя")
	}

	return c.JSON(user)
}

func (a *AuthControllers) UserSignUp(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	c.BodyParser(&credentials)

	newUser := models.User{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Login:     credentials.Login,
		Email:     credentials.Email,
		Password:  credentials.Password,
		Gender:    credentials.Gender,
		Birthday:  credentials.Birthday,
		Location:  credentials.Location,
		Job:       credentials.Job,
	}

	newUserObj, err := a.UserServices.InsertOne(newUser)
	if err != nil {
		return err
	}

	return c.JSON(newUserObj)
}

func (a *AuthControllers) ValidateCredentials(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	var isCredentialsValid models.IsCredentials
	c.BodyParser(&credentials)

	isValid, _ := a.UserServices.FindForValidate(credentials.Login, credentials.Email)

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
		unvalidMsg, _ := json.Marshal(&isCredentialsValid)
		return fiber.NewError(409, string(unvalidMsg))
	}

	return c.JSON(isCredentialsValid)
}

func GetAuthControllers(store store.InterfaceStore) AuthControllers {
	return AuthControllers{
		UserServices: store.Users(),
	}
}
