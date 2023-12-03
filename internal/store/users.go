package store

import (
	"fmt"
	"github.com/Troom-Corp/troom/internal/models"
	"github.com/Troom-Corp/troom/internal/pkg"
	"github.com/jmoiron/sqlx"
)

type InterfaceUser interface {
	InsertOne(user models.User) (int, error)
	DeleteOne(userid int) error
	FindByQuery(searchQuery string) ([]models.User, error)
	FindByLogin(login string) (models.User, error)
	UserExists(login string) (models.User, error)
	UpdateOne(user models.User) error
	FindForValidate(login, email string) ([]models.User, error)
	FindByID(userID int) (models.User, error)
	UploadAvatar(userID int, filename string) (string, error)
	UploadLayout(userID int, filename string) (string, error)
	DeleteAvatar(userID int) (string, error)
	DeleteLayout(userID int) (string, error)
}

type user struct {
	db *sqlx.DB
}

func (u user) InsertOne(user models.User) (int, error) {
	var insertedID int

	passwordHash, err := pkg.Encode([]byte(user.Password))
	if err != nil {
		return insertedID, err
	}

	err = u.db.Get(&insertedID, fmt.Sprintf("insert into users (role, firstname, lastname, login, email, password, gender, birthday, location, job, phone, links, avatar, bio, layout) "+
		"values ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s') RETURNING userid",
		"users", user.FirstName, user.LastName, user.Login, user.Email, passwordHash, user.Gender, user.Birthday, user.Location, user.Job, user.Phone, user.Links, user.Avatar, user.Bio, user.Layout))

	return insertedID, err
}

func (u user) DeleteOne(userid int) error {
	_, err := u.db.Query(fmt.Sprintf("delete from users where userid = %d", userid))
	return err
}

func (u user) FindByQuery(searchQuery string) ([]models.User, error) {
	var queryUsers []models.User

	searchFormat := "%" + searchQuery + "%"
	err := u.db.Select(&queryUsers, fmt.Sprintf("select * from users where lower(firstname) like lower('%s') or lower(lastname) like lower('%s') or lower(login) like lower('%s') LIMIT 5", searchFormat, searchFormat, searchFormat))

	return queryUsers, err
}

func (u user) FindByLogin(login string) (models.User, error) {
	var resultUser models.User

	err := u.db.Get(&resultUser, fmt.Sprintf("select * from users where login = '%s'", login))

	return resultUser, err
}

func (u user) UserExists(login string) (models.User, error) {
	var resultUser models.User

	err := u.db.Get(&resultUser, fmt.Sprintf("select * from users where login = '%s' or email = '%s'", login, login))

	return resultUser, err
}

func (u user) FindForValidate(login, email string) ([]models.User, error) {
	var resultUsers []models.User

	err := u.db.Select(&resultUsers, fmt.Sprintf("select * from users where login = '%s' union select * from users where email = '%s'", login, email))

	return resultUsers, err
}

func (u user) UpdateOne(user models.User) error {
	passwordHash, err := pkg.Encode([]byte(user.Password))
	if err != nil {
		return err
	}

	_, err = u.db.Query(fmt.Sprintf("update users set "+
		"firstname='%s', secondname='%s', login='%s', email='%s', password='%s', gender='%s', birthday='%s', location='%s', job= %s', phone='%s', links='%s', avatar='%s', bio='%s' where userid = %d",
		user.FirstName, user.LastName, user.Login, user.Email, passwordHash, user.Gender, user.Birthday, user.Location, user.Job, user.Phone, user.Links, user.Avatar, user.Bio, user.UserId))

	return err
}

func (u user) FindByID(userID int) (models.User, error) {
	var profile models.User

	err := u.db.Get(&profile, fmt.Sprintf("select * from users where userid = %d", userID))

	return profile, err
}

func (u user) UploadAvatar(userID int, filename string) (string, error) {
	var oldAvatar string
	err := u.db.Get(&oldAvatar, fmt.Sprintf("update users set avatar = '%s' where userid = %d returning (select avatar from users where userid = %d)", filename, userID, userID))

	return oldAvatar, err
}

func (u user) UploadLayout(userID int, filename string) (string, error) {
	var oldLayout string
	err := u.db.Get(&oldLayout, fmt.Sprintf("update users set layout = '%s' where userid = %d returning (select layout from users where userid = %d)", filename, userID, userID))

	return oldLayout, err

}

func (u user) DeleteAvatar(userID int) (string, error) {
	var deletedAvatar string
	err := u.db.Get(&deletedAvatar, fmt.Sprintf("update users set avatar = '' where userid = %d returning (select avatar from users where userid = %d)", userID, userID))
	return deletedAvatar, err
}

func (u user) DeleteLayout(userID int) (string, error) {
	var deletedLayout string
	err := u.db.Get(&deletedLayout, fmt.Sprintf("update users set layout = '' where userid = %d returning (select layout from users where userid = %d)", userID, userID))
	return deletedLayout, err
}
