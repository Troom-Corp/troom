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
	Update() error
	Delete() error
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

	conn.Close()
	return queryUsers, nil
}

// Update Обновить данные пользователя по ID
func (u User) Update() error {
	var userId int
	var userEmail, userNick, userPhone string

	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	updateByIdQuery := fmt.Sprintf("UPDATE public.users SET "+
		"firstname = '%s', secondname = '%s', email = '%s', avatar = '%s', bio = '%s', "+
		"phone = '%s', links = '%s', nick = '%s' WHERE userid = %d IF email NOT ",
		u.FirstName, u.SecondName, u.Email, u.Avatar, u.Bio, u.Phone, u.Links, u.Nick, u.UserId)
	duplicateEmailAndNick := fmt.Sprintf("SELECT email, nick, phone FROM public.users WHERE NOT userid = %d AND (email = '%s' OR nick = '%s' OR phone = '%s')", u.UserId, u.Email, u.Phone, u.Nick)
	duplicateId := fmt.Sprintf("SELECT userid FROM public.users WHERE userid = %d", u.UserId)

	conn.Get(&userId, duplicateId)
	if userId == 0 {
		conn.Close()
		return fiber.NewError(404, "Такого пользователя не сущесвует")
	}

	rows, _ := conn.Queryx(duplicateEmailAndNick)
	for rows.Next() {
		rows.Scan(&userEmail, &userNick, &userPhone)
	}
	if userEmail == u.Email {
		conn.Close()
		return fiber.NewError(409, "Такая почта уже существует")
	}
	if userNick == u.Nick {
		conn.Close()
		return fiber.NewError(409, "Такой ник уже существует")
	}
	if userPhone == u.Phone {
		conn.Close()
		return fiber.NewError(409, "Такой номер телефона уже существует")
	}

	_, err = conn.Query(updateByIdQuery)
	conn.Close()
	return err
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

	err = conn.Get(&userId, deleteByIdQuery)
	if userId == 0 {
		conn.Close()
		return fiber.NewError(409, "Пользователя не сущесвует")
	}

	conn.Close()
	return err
}
