package services

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/storage"
)

type UserInterface interface {
	Create() (int, error)
	ReadAll() ([]User, error)
	ReadById() (User, error)
	SearchByQuery(string) ([]User, error)
	Update() error
	Delete() error
}

type User struct {
	UserId      int    `db:"userid"`
	Password    string `db:"password"`
	Role        string `db:"role"`
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
		return 0, err
	}

	createQuery := fmt.Sprintf("INSERT INTO "+
		"public.users (password, role, firstname, secondname, gender, age, dateofbirth, location, phone, email, links, job, avatar, bio) "+
		"VALUES ('%s', 'user', '%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid",
		u.Password, u.FirstName, u.SecondName, u.Gender, u.Age, u.DateOfBirth, u.Location, u.Phone, u.Email, u.Links, u.Job, u.Avatar, u.Bio)

	conn.Get(&userId, createQuery)
	err = conn.Close()
	return userId, err
}

// ReadAll Прочитать всех пользователей
func (u User) ReadAll() ([]User, error) {
	var users []User
	conn, err := storage.Sql.Open()

	if err != nil {
		return []User{}, err
	}

	err = conn.Select(&users, "SELECT * FROM public.users")

	if err != nil {
		fmt.Println(err)
		conn.Close()
		return []User{}, err
	}

	conn.Close()
	return users, nil
}

// ReadById Прочитать одного пользователя по ID
func (u User) ReadById() (User, error) {
	var user User
	conn, err := storage.Sql.Open()

	if err != nil {
		return User{}, err
	}

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.users WHERE userid=%d", u.UserId)
	err = conn.Get(&user, readByIdQuery)

	if err != nil {
		conn.Close()
		return User{}, err
	}

	conn.Close()
	return user, nil
}

func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	conn, err := storage.Sql.Open()

	if err != nil {
		return []User{}, err
	}

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.users WHERE LOWER(firstname) LIKE '%s' OR LOWER(secondname) LIKE '%s'", searchFormat, searchFormat)
	err = conn.Select(&queryUsers, searchByQuery)

	if err != nil {
		conn.Close()
		return []User{}, nil
	}

	conn.Close()
	return queryUsers, nil
}

// Update Обновить данные пользователя по ID
func (u User) Update() error {
	conn, err := storage.Sql.Open()
	if err != nil {
		return err
	}

	updateByIdQuery := fmt.Sprintf("UPDATE public.users SET firstname = '%s', secondname = '%s', email = '%s', password = '%s', photo = '%s', bio = '%s', phone = '%s', links = '%s', followers = '%s', subscribers = '%s' WHERE userid = %d", u.FirstName, u.SecondName, u.Email, u.Password, u.Avatar, u.Bio, u.Phone, u.Links, u.Role, u.DateOfBirth, u.UserId)
	_, err = conn.Query(updateByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}

// Delete Удалить все данные пользователя по ID
func (u User) Delete() error {
	conn, err := storage.Sql.Open()

	if err != nil {
		return err
	}

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.users WHERE userid = %d", u.UserId)
	_, err = conn.Query(deleteByIdQuery)

	if err != nil {
		conn.Close()
		return err
	}

	conn.Close()
	return nil
}
