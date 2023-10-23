package services

import (
	"context"
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
	UserId      int
	FirstName   string
	SecondName  string
	Email       string
	Password    string
	Photo       string
	Bio         string
	Phone       string
	Links       string
	Followers   string
	Subscribers string
}

// Create Создать пользователя по входным данным и получить ID этого пользователя
func (u User) Create() (int, error) {
	var userId int
	conn := storage.SqlInterface.New()

	createQuery := fmt.Sprintf("INSERT INTO public.users (firstname, secondname, email, password, photo, bio, phone, links, followers, subscribers) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid", u.FirstName, u.SecondName, u.Email, u.Password, u.Photo, u.Bio, u.Phone, u.Links, u.Followers, u.Subscribers)
	rows, err := conn.Query(context.Background(), createQuery)
	rows.Scan(&userId)

	storage.SqlInterface.Close(conn)
	return userId, err
}

// ReadAll Прочитать всех пользователей
func (u User) ReadAll() ([]User, error) {
	var users []User
	conn := storage.SqlInterface.New()

	rows, err := conn.Query(context.Background(), "SELECT * FROM public.users;")

	if err != nil {
		storage.SqlInterface.Close(conn)
		return []User{}, err
	}

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.UserId,
			&user.FirstName,
			&user.SecondName,
			&user.Email,
			&user.Password,
			&user.Photo,
			&user.Bio,
			&user.Phone,
			&user.Links,
			&user.Followers,
			&user.Subscribers)
		if err != nil {
			storage.SqlInterface.Close(conn)
			return []User{}, err
		}
		users = append(users, user)
	}

	storage.SqlInterface.Close(conn)
	return users, nil
}

// ReadById Прочитать одного пользователя по ID
func (u User) ReadById() (User, error) {
	var user User
	conn := storage.SqlInterface.New()

	readByIdQuery := fmt.Sprintf("SELECT * FROM public.users WHERE userid=%d", u.UserId)
	err := conn.QueryRow(context.Background(), readByIdQuery).Scan(&user.UserId, &user.FirstName, &user.SecondName, &user.Email, &user.Password, &user.Photo, &user.Bio, &user.Phone, &user.Links, &user.Followers, &user.Subscribers)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return User{}, err
	}

	storage.SqlInterface.Close(conn)
	return user, nil
}

func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	conn := storage.SqlInterface.New()

	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.users WHERE LOWER(firstname) LIKE '%s' OR LOWER(secondname) LIKE '%s'", searchFormat, searchFormat)
	rows, err := conn.Query(context.Background(), searchByQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return []User{}, nil
	}

	for rows.Next() {
		var queryUser User
		err = rows.Scan(
			&queryUser.UserId,
			&queryUser.FirstName,
			&queryUser.SecondName,
			&queryUser.Email,
			&queryUser.Password,
			&queryUser.Photo,
			&queryUser.Bio,
			&queryUser.Phone,
			&queryUser.Links,
			&queryUser.Followers,
			&queryUser.Subscribers)
		if err != nil {
			storage.SqlInterface.Close(conn)
			return []User{}, err
		}
		queryUsers = append(queryUsers, queryUser)
	}

	storage.SqlInterface.Close(conn)
	return queryUsers, nil
}

// Update Обновить данные пользователя по ID
func (u User) Update() error {
	conn := storage.SqlInterface.New()

	updateByIdQuery := fmt.Sprintf("UPDATE public.users SET firstname = '%s', secondname = '%s', email = '%s', password = '%s', photo = '%s', bio = '%s', phone = '%s', links = '%s', followers = '%s', subscribers = '%s' WHERE userid = %d", u.FirstName, u.SecondName, u.Email, u.Password, u.Photo, u.Bio, u.Phone, u.Links, u.Followers, u.Subscribers, u.UserId)
	_, err := conn.Query(context.Background(), updateByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}

// Delete Удалить все данные пользователя по ID
func (u User) Delete() error {
	conn := storage.SqlInterface.New()

	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.users WHERE userid = %d", u.UserId)
	_, err := conn.Query(context.Background(), deleteByIdQuery)

	if err != nil {
		storage.SqlInterface.Close(conn)
		return err
	}

	storage.SqlInterface.Close(conn)
	return nil
}
