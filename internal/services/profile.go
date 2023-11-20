package services

import (
	"encoding/json"
	"fmt"
	"github.com/Troom-Corp/troom/internal/broker"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type UserIdentification struct {
	UserId int    `json:"userid"`
	Email  string `json:"email"`
}

type ProfileInfo struct {
	UserId      int    `json:"-"`
	FirstName   string `json:"firstName"`
	SecondName  string `json:"secondName"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Location    string `json:"location"`
	Job         string `json:"job"`
	Links       string `json:"links"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
}

func (pi ProfileInfo) UserProfile() (User, error) {
	var userProfile User
	conn, err := storage.Sql.Open()
	if err != nil {
		return User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}
	getProfileQuery := fmt.Sprintf("SELECT * FROM users WHERE userid = %d", pi.UserId)
	err = conn.Get(&userProfile, getProfileQuery)

	if err != nil {
		conn.Close()
		return User{}, fiber.NewError(404, "Ошибка при открытии профиля")
	}

	conn.Close()
	return userProfile, nil
}

func (pi ProfileInfo) UpdateInfo() error {
	var userId int

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	updateInfoQuery := fmt.Sprintf("UPDATE public.users SET "+
		"firstname = '%s', secondname = '%s', gender = '%s', dateofbirth = '%s', location = '%s', job = '%s', links = '%s', avatar = '%s', bio = '%s' WHERE userid = %d RETURNING userid",
		pi.FirstName, pi.SecondName, pi.Gender, pi.DateOfBirth, pi.Location, pi.Job, pi.Links, pi.Avatar, pi.Bio, pi.UserId)
	err = conn.Get(&userId, updateInfoQuery)

	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не существует")
	}
	err = conn.Close()
	return err
}

type NewPasswordCredentials struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (np NewPasswordCredentials) GetResetPasswordLink(userid int) error {
	var userIdentity UserIdentification
	userIdentity.UserId = userid

	var userId int
	var userPassword string

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getUserQuery := fmt.Sprintf("SELECT userid, password, email FROM public.users WHERE userid = %d", userid)

	rows, err := conn.Query(getUserQuery)

	for rows.Next() {
		rows.Scan(&userId, &userPassword, &userIdentity.Email)
	}

	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не существует")
	}

	if pkg.Decode([]byte(userPassword), []byte(np.OldPassword)) != nil {
		conn.Close()
		return fiber.NewError(401, "Неправильный пароль")
	}

	newPasswordData, _ := json.Marshal(userIdentity)

	rabbitConn, _ := broker.RabbitBroker.Connect()
	broker.Rabbit{Broker: rabbitConn}.SendMsg("resetPassword", newPasswordData)

	broker.Rabbit{Broker: rabbitConn}.Close()
	conn.Close()
	return fiber.NewError(200, "Ссылка была отправлена на вашу почту")
}

func (np NewPasswordCredentials) SetNewPassword(userid int) error {

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}
	hashedPassword, _ := pkg.Encode([]byte(np.NewPassword))
	setNewPasswordQuery := fmt.Sprintf("UPDATE public.users SET password = '%s' WHERE userid = %d", hashedPassword, userid)

	_, err = conn.Query(setNewPasswordQuery)
	conn.Close()
	return err
}

type NewEmailCredentials struct {
	NewEmail string `json:"new_email"`
}

func (ne NewEmailCredentials) GetResetEmailLink(userid int) error {
	var userIdentity UserIdentification
	userIdentity.UserId = userid

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}
	defer conn.Close()
	getUserEmailQuery := fmt.Sprintf("SELECT email FROM public.users WHERE userid = %d", userid)
	conn.Get(&userIdentity.Email, getUserEmailQuery)

	if userIdentity.Email == "" {
		return fiber.NewError(404, "Пользователя не сущесвует")
	}

	newEmailData, _ := json.Marshal(userIdentity)

	rabbitConn, _ := broker.RabbitBroker.Connect()
	broker.Rabbit{Broker: rabbitConn}.SendMsg("resetEmail", newEmailData)

	broker.Rabbit{Broker: rabbitConn}.Close()
	conn.Close()
	return fiber.NewError(200, "Ссылка была отправлена на вашу почту")
}

func (ne NewEmailCredentials) SetNewEmail(userid int) error {
	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	setNewPasswordQuery := fmt.Sprintf("UPDATE public.users SET email = '%s' WHERE userid = %d", ne.NewEmail, userid)
	fmt.Println(setNewPasswordQuery)
	_, err = conn.Query(setNewPasswordQuery)
	conn.Close()
	return err
}

type NewLoginCredentials struct {
	NewLogin string `json:"new_login"`
}

func (nl NewLoginCredentials) SetNewLogin(userid int) error {
	var userLogin string

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getLoginQuery := fmt.Sprintf("SELECT nick FROM public.users WHERE nick = '%s'", nl.NewLogin)
	updateLoginQuery := fmt.Sprintf("UPDATE public.users SET nick = '%s' WHERE userid = %d RETURNING nick", nl.NewLogin, userid)
	conn.Get(&userLogin, getLoginQuery)

	if userLogin != "" || nl.NewLogin == "" {
		conn.Close()
		return fiber.NewError(409, "Такой логин уже используется")
	}

	conn.Query(updateLoginQuery)
	conn.Close()
	return fiber.NewError(200, "Логин успешно обновлен")
}

type NewPhoneCredentials struct {
	OldPhoneCredentials string
	VerificationSMS     string
	NewPhone            string
}

func (np NewPhoneCredentials) UpdatePhone() error {
	return nil
}
