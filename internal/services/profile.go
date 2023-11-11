package services

import (
	"encoding/json"
	"fmt"
	"github.com/Troom-Corp/troom/internal/broker"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type ProfileInterface interface {
	UpdateInfo(userid int) error // FirstName, SecondName, Gender, DateOfBirth, Location, Job, Links, Avatar, Bio

	// need to check for duplicates
	ResetPassword(userid int) error
	CheckCode(userid int) error
	UpdatePhone() error
	UpdateEmail() error
	UpdateLogin(newLogin string, userid int) error
}

type ProfileInfo struct {
	FirstName   string
	SecondName  string
	Gender      string
	DateOfBirth string
	Location    string
	Job         string
	Links       string
	Avatar      string
	Bio         string
}

type UserIdentification struct {
	UserId int    `json:"userid"`
	Email  string `json:"email"`
}

type NewPasswordCredentials struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (np NewPasswordCredentials) GetResetLink(userid int) error {
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
	broker.Rabbit{Broker: rabbitConn}.SendMsg("verifyPassword", newPasswordData)

	broker.Rabbit{Broker: rabbitConn}.Close()
	conn.Close()
	return nil
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
	OldEmail         string
	VerificationCode string
	NewEmail         string
}

type NewPhoneCredentials struct {
	OldPhoneCredentials string
	VerificationSMS     string
	NewPhone            string
}

// UpdateInfo Обновляет данные профиля пользователя
func (p ProfileInfo) UpdateInfo(userid int) error {
	var userId int

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	updateInfoQuery := fmt.Sprintf("UPDATE public.users SET "+
		"firstname = '%s', secondname = '%s', gender = '%s', dateofbirth = '%s', location = '%s', job = '%s', links = '%s', avatar = '%s', bio = '%s' WHERE userid = %d RETURNING userid",
		p.FirstName, p.SecondName, p.Gender, p.DateOfBirth, p.Location, p.Job, p.Links, p.Avatar, p.Bio, userid)
	err = conn.Get(&userId, updateInfoQuery)

	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не существует")
	}
	err = conn.Close()
	return err
}

func (np NewPhoneCredentials) UpdatePhone() error {
	return nil
}

func (ne NewEmailCredentials) UpdateEmail() error {
	var userEmail string

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getEmailQuery := fmt.Sprintf("SELECT email FROM public.users WHERE email = '%s' RETURNING email", ne.OldEmail)
	updateEmailQuery := fmt.Sprintf("UPDATE SET email = '%s'", ne.NewEmail)
	conn.Get(&userEmail, getEmailQuery)

	if userEmail != "" {
		conn.Close()
		return fiber.NewError(409, "Такая почта уже используется")
	}

	// здесь должна быть логика брокера сообщений

	conn.Query(updateEmailQuery)
	conn.Close()
	return nil
}

func (p ProfileInfo) UpdateLogin(newLogin string, userid int) error {
	var userLogin string

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getLoginQuery := fmt.Sprintf("SELECT nick FROM public.users WHERE nick = '%s'", newLogin)
	updateLoginQuery := fmt.Sprintf("UPDATE public.users SET nick = '%s' WHERE userid = %d RETURNING nick", newLogin, userid)
	conn.Get(&userLogin, getLoginQuery)

	if userLogin != "" {
		conn.Close()
		return fiber.NewError(409, "Такой nick уже используется")
	}

	conn.Query(updateLoginQuery)
	conn.Close()
	return nil
}
