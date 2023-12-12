package controllers

import (
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/store"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type AuthControllers struct {
	UserServices store.InterfaceUser
}

func (a AuthControllers) SignIn(c *fiber.Ctx) error {
	var credentials models.SignInCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad JSON form",
			},
		})
	}

	user, _ := a.UserServices.IsUserExist(credentials.Login)

	if err := pkg.Decode([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return c.Status(401).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "401",
				Message: "Invalid credentials",
			},
		})
	}

	token, err := pkg.SignJWT(user.UserId)
	if err != nil {
		return c.Status(401).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "401",
				Message: "An error while creating an user token",
			},
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: "None",
		Secure:   false,
	})

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "You successfully logged in",
		},
	})
}

func (a AuthControllers) SignUp(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	err := c.BodyParser(&credentials)

	if err != nil {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad JSON form",
			},
		})
	}

	newUserCredentials := models.User{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Login:     strings.ToLower(credentials.Login),
		Email:     credentials.Email,
		Password:  credentials.Password,
	}

	if !credentials.Validate() {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad credentials",
			},
		})
	}

	insertedID, err := a.UserServices.InsertOne(newUserCredentials)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while creating an user",
			},
		})
	}

	token, err := pkg.SignJWT(insertedID)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while creating an user token",
			},
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: "None",
		Secure:   false,
	})

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "You successfully create an account",
		},
	})
}

func (a AuthControllers) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
	})

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "You successfully logged out",
		},
	})
}

// ValidateCredentials runs on the client before the SignUp method
func (a AuthControllers) ValidateCredentials(c *fiber.Ctx) error {
	var credentials models.SignUpCredentials
	var isCredentialsValid models.IsCredentials
	err := c.BodyParser(&credentials)
	if err != nil {
		return c.Status(400).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "400",
				Message: "Bad JSON form",
			},
		})
	}

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
		return c.Status(409).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "409",
				Message: "Bad Credentials",
			},
			Data: isCredentialsValid,
		})
	}

	return c.JSON(isCredentialsValid)
}

func (a AuthControllers) Profile(c *fiber.Ctx) error {
	jwt := c.Cookies("token")
	ID, err := pkg.GetIdentity(jwt)
	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while identification user",
			},
		})
	}

	userProfile, err := a.UserServices.FindOne("userid", ID)

	if err != nil {
		return c.Status(500).JSON(models.HttpResponse{
			Error: models.Error{
				Status:  "500",
				Message: "An error while getting an user",
			},
		})
	}

	return c.JSON(models.HttpResponse{
		Error: models.Error{
			Status:  "200",
			Message: "User",
		},
		Data: userProfile,
	})
}

func GetAuthControllers(store store.InterfaceStore) *AuthControllers {
	return &AuthControllers{
		UserServices: store.Users(),
	}
}
