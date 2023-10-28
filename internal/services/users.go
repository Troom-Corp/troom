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
	UserId      int    `db:"userid"`
	Password    string `db:"password"`
	Role        string `db:"role"`
	Nick        string `db:"nick"`
	FirstName   string `db:"firstname"`
	SecondName  string `db:"secondname"`
	Gender      string `db:"gender"`
	Age         int    `db:"age"`
	DateOfBirth string `db:"dateofbirth"`
	Location    string `db:"location"`
	Phone       string `db:"phone"`
	Email       string `db:"email" `
	Links       string `db:"links"`
	Job         string `db:"job"`
	Avatar      string `db:"avatar"`
	Bio         string `db:"bio"`
}

// Create Создать пользователя по входным данным и получить ID этого пользователя
func (u User) Create() (int, error) {
	var userId int
	conn, err := storage.Sql.Open()

	if err != nil {
		return 0, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	createQuery := fmt.Sprintf("INSERT INTO "+
		"public.users (password, role, firstname, secondname, gender, age, dateofbirth, location, phone, email, links, job, avatar, bio, nick) "+
		"VALUES ('%s', 'user', '%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid",
		u.Password, u.FirstName, u.SecondName, u.Gender, u.Age, u.DateOfBirth, u.Location, u.Phone, u.Email, u.Links, u.Job, u.Avatar, u.Bio, u.Nick)
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
		return []User{}, fiber.NewError(500, "Ошибка при получении пользователей")
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

func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	conn, err := storage.Sql.Open()

	if err != nil {
		return []User{}, fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

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
	conn, err := storage.Sql.Open()
	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	updateByIdQuery := fmt.Sprintf("UPDATE public.users SET "+
		"firstname = '%s', secondname = '%s', email = '%s', avatar = '%s', bio = '%s', "+
		"phone = '%s', links = '%s' WHERE userid = %d",
		u.FirstName, u.SecondName, u.Email, u.Avatar, u.Bio, u.Phone, u.Links, u.UserId)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return fiber.NewError(500, "Ошибка при обновлении профиля")
	}

	conn.Close()
	return nil
}

// Delete Удалить все данные пользователя по ID
func (u User) Delete() error {
	conn, err := storage.Sql.Open()

	if err != nil {
		return fiber.NewError(500, "Ошибка при подключении к базе данных")
	}

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.users WHERE userid = %d", u.UserId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return fiber.NewError(500, "Ошибка при удалении профиля")
	}

	conn.Close()
	return nil
}
