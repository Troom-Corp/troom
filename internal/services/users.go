package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type UserInterface interface {
	Create() (int, error)
	ReadAll() ([]User, error)
	ReadByNick() (User, error)
	SearchByQuery(string) ([]User, error)
	Delete() error
}

type ProfileInterface interface {
	UpdateInfo() error // FirstName, SecondName, Gender, DateOfBirth, Location, Job, Links, Avatar, Bio

	// need to check for duplicates
	ResetPassword() error
	UpdatePhone() error
	UpdateEmail() error
	UpdateNick() error
}

type User struct {
	UserId int    `db:"userid"`
	Role   string `db:"role"`

	// первый этап регистрации
	FirstName  string `db:"firstname"`
	SecondName string `db:"secondname"`
	Nick       string `db:"nick"`
	Email      string `db:"email" `
	Password   string `db:"password"`

	// второй этап регистрации
	Gender      string `db:"gender"`
	DateOfBirth string `db:"dateofbirth"`
	Location    string `db:"location"`
	Job         string `db:"job"`

	// настраивается в профиле
	Phone  string `db:"phone"`  // in profile
	Links  string `db:"links"`  // in profile
	Avatar string `db:"avatar"` // in profile
	Bio    string `db:"bio"`    // in profile
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
		u.FirstName, u.SecondName, u.Nick, u.Email, u.Password, u.Gender, u.DateOfBirth, u.Location, u.Job, u.Phone, u.Links, u.Avatar, u.Bio)

	err = conn.Get(&userId, createQuery)
	if err != nil {
		fmt.Println(err)
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

// ReadByNick Прочитать одного пользователя по ID
func (u User) ReadByNick() (User, error) {
	var user User
	conn, err := storage.Sql.Open()

	if err != nil {
		return User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.users WHERE nick='%s'", u.Nick)

	err = conn.Get(&user, readByIdQuery)
	if err != nil {
		conn.Close()
		return User{}, fiber.NewError(404, "Пользователь не найден")
	}

	conn.Close()
	return user, nil
}

// SearchByQuery Найти пользователей по searchQuery
func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	conn, err := storage.Sql.Open()

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM users WHERE LOWER(firstname) LIKE LOWER('%s') OR LOWER(secondname) LIKE LOWER('%s')", searchFormat, searchFormat)

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
		conn.Close()
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

// UpdateInfo Обновляет данные профиля пользователя
func (u User) UpdateInfo() error {
	var userId int

	conn, err := storage.Sql.Open()
	if err != nil {
		conn.Close()
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	updateInfoQuery := fmt.Sprintf("UPDATE SET "+
		"firstname = '%s', secondname = '%s', gender = '%s', dateofbirth = '%s', location = '%s', job = '%s', links = '%s', avatar = '%s', bio = '%s' WHERE userid = %d RETURNING userid",
		u.FirstName, u.SecondName, u.Gender, u.DateOfBirth, u.Location, u.Job, u.Links, u.Avatar, u.Bio, u.UserId)
	conn.Get(&userId, updateInfoQuery)
	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не сущесвует")
	}
	err = conn.Close()
	return err
}

func (u User) ResetPassword() error {
	return nil
}

func (u User) UpdatePhone() error {
	return nil
}

func (u User) UpdateEmail() error {
	var userEmail string

	conn, err := storage.Sql.Open()
	if err != nil {
		conn.Close()
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getEmailQuery := fmt.Sprintf("SELECT email FROM public.users WHERE email = '%s' RETURNING email", u.Email)
	updateEmailQuery := fmt.Sprintf("UPDATE SET email = '%s'", u.Email)
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

func (u User) UpdateNick() error {
	var userNick string

	conn, err := storage.Sql.Open()
	if err != nil {
		conn.Close()
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	getNickQuery := fmt.Sprintf("SELECT nick FROM public.users WHERE nick = '%s' RETURNING nick", u.Nick)
	updateNickQuery := fmt.Sprintf("UPDATE SET nick = '%s'", u.Nick)
	conn.Get(&userNick, getNickQuery)

	if userNick != "" {
		conn.Close()
		return fiber.NewError(409, "Такой nick уже используется")
	}

	conn.Query(updateNickQuery)
	conn.Close()
	return nil
}
