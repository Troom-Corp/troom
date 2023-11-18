package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type UserInterface interface {
	Create() (int, error)
	ReadAll() ([]User, error)
	ReadByLogin() (User, error)
	SearchByQuery(string) ([]User, error)
	UserProfile() (User, error)
	Delete() error
}

type User struct {
	UserId int    `db:"userid" json:"id"`
	Role   string `db:"role" json:"-"`

	// первый этап регистрации
	FirstName  string `db:"firstname" json:"firstName"`
	SecondName string `db:"secondname" json:"secondName"`
	Login      string `db:"nick" json:"login"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password" json:"-"`

	// второй этап регистрации
	Gender      string `db:"gender" json:"gender"`
	DateOfBirth string `db:"dateofbirth" json:"date_of_birth"`
	Location    string `db:"location" json:"location"`
	Job         string `db:"job" json:"job"`

	// настраивается в профиле
	Phone  string `db:"phone" json:"phone"`   // in profile
	Links  string `db:"links" json:"links"`   // in profile
	Avatar string `db:"avatar" json:"avatar"` // in profile
	Bio    string `db:"bio" json:"bio"`       // in profile
}

// Create Создать пользователя по входным данным и получить ID этого пользователя
func (u User) Create() (int, error) {
	var userId int
	conn, err := storage.Sql.Open()

	if err != nil {
		return 0, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	createQuery := fmt.Sprintf("INSERT INTO "+
		"public.users (role, firstname, secondname, nick, email, password, gender, dateofbirth, location, job, phone, links, avatar, bio) "+
		"VALUES ('user', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid",
		u.FirstName, u.SecondName, u.Login, u.Email, u.Password, u.Gender, u.DateOfBirth, u.Location, u.Job, u.Phone, u.Links, u.Avatar, u.Bio)

	err = conn.Get(&userId, createQuery)
	if err != nil {
		return 0, fiber.NewError(500, "Ошибка при создании пользователя")
	}

	conn.Close()
	return userId, nil
}

// ReadAll Прочитать всех пользователей
func (u User) ReadAll() ([]User, error) {
	var users []User
	conn, err := storage.Sql.Open()

	if err != nil {
		return []User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	err = conn.Select(&users, "SELECT * FROM public.users")
	if err != nil {
		conn.Close()
		return []User{}, fiber.NewError(500, "Неизвестная ошибка")
	}

	conn.Close()
	return users, nil
}

// ReadByLogin Прочитать одного пользователя по ID
func (u User) ReadByLogin() (User, error) {
	var user User
	conn, err := storage.Sql.Open()

	if err != nil {
		return User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.users WHERE nick='%s'", u.Login)

	err = conn.Get(&user, readByIdQuery)
	if err != nil {
		conn.Close()
		return User{}, fiber.NewError(404, "Пользователь не найден")
	}

	conn.Close()
	return user, nil
}

func (u User) UserProfile() (User, error) {
	var userProfile User
	conn, err := storage.Sql.Open()
	if err != nil {
		return User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}
	getProfileQuery := fmt.Sprintf("SELECT * FROM users WHERE userid = %d", u.UserId)
	err = conn.Get(&userProfile, getProfileQuery)

	if err != nil {
		conn.Close()
		return User{}, fiber.NewError(404, "Ошибка при открытии профиля")
	}

	conn.Close()
	return userProfile, nil
}

// SearchByQuery Найти пользователей по searchQuery
func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	conn, err := storage.Sql.Open()

	if err != nil {
		return []User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM users WHERE LOWER(firstname) LIKE LOWER('%s') OR LOWER(secondname) LIKE LOWER('%s') OR LOWER(nick) LIKE LOWER('%s') LIMIT 5", searchFormat, searchFormat, searchFormat)

	err = conn.Select(&queryUsers, searchByQuery)
	if err != nil {
		conn.Close()
		return []User{}, fiber.NewError(500, "Ошибка при поиске пользователей")
	}

	err = conn.Close()
	return queryUsers, err
}

// Delete Удалить все данные пользователя по ID
func (u User) Delete() error {
	var userId int

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.users WHERE userid = %d RETURNING userid", u.UserId)

	conn.Get(&userId, deleteByIdQuery)
	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не сущесвует")
	}

	err = conn.Close()
	return err
}
