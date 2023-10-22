package services

import (
	"context"
	"fmt"

	"github.com/Troom-Corp/troom/internal"
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
	createQuery := fmt.Sprintf("INSERT INTO public.users (firstname, secondname, email, password, photo, bio, phone, links, followers, subscribers) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid",
		u.FirstName,
		u.SecondName,
		u.Email,
		u.Password,
		u.Photo,
		u.Bio,
		u.Phone,
		u.Links,
		u.Followers,
		u.Subscribers)
	rows, err := internal.Store().Query(context.Background(), createQuery)
	rows.Scan(&userId)
	return userId, err
}

// ReadAll Прочитать всех пользователей
func (u User) ReadAll() ([]User, error) {
	var users []User

	rows, _ := internal.Store().Query(context.Background(), "SELECT * FROM public.users;")
	for rows.Next() {
		var user User
		err := rows.Scan(
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
			return []User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

// ReadById Прочитать одного пользователя по ID
func (u User) ReadById() (User, error) {
	var user User
	readByIdQuery := fmt.Sprintf("SELECT * FROM public.users WHERE userid=%d", u.UserId)
	err := internal.Store().QueryRow(context.Background(), readByIdQuery).Scan(
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
		return User{}, err
	}
	return user, nil
}

func (u User) SearchByQuery(searchQuery string) ([]User, error) {
	var queryUsers []User
	searchFormat := "%" + searchQuery + "%"
	searchByQuery := fmt.Sprintf("SELECT * FROM public.users WHERE LOWER(firstname) LIKE '%s' OR LOWER(secondname) LIKE '%s'", searchFormat, searchFormat)
	rows, err := internal.Store().Query(context.Background(), searchByQuery)
	if err != nil {
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
			return []User{}, err
		}
		queryUsers = append(queryUsers, queryUser)
	}

	return queryUsers, nil
}

// Update Обновить данные пользователя по ID
func (u User) Update() error {
	updateByIdQuery := fmt.Sprintf("UPDATE public.users SET firstname = '%s', secondname = '%s', email = '%s', password = '%s', photo = '%s', bio = '%s', phone = '%s', links = '%s', followers = '%s', subscribers = '%s' WHERE userid = %d",
		u.FirstName,
		u.SecondName,
		u.Email,
		u.Password,
		u.Photo,
		u.Bio,
		u.Phone,
		u.Links,
		u.Followers,
		u.Subscribers,
		u.UserId)
	_, err := internal.Store().Query(context.Background(), updateByIdQuery)
	if err != nil {
		return err
	}
	return nil
}

// Delete Удалить все данные пользователя по ID
func (u User) Delete() error {
	deleteByIdQuery := fmt.Sprintf("DELETE FROM public.users WHERE userid = %d", u.UserId)
	_, err := internal.Store().Query(context.Background(), deleteByIdQuery)
	if err != nil {
		return err
	}
	return nil
}
