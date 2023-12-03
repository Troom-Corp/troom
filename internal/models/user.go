package models

type User struct {
	UserId int    `db:"userid" json:"id"`
	Role   string `db:"role" json:"-"`

	// первый этап регистрации
	FirstName string `db:"firstname" json:"firstName"`
	LastName  string `db:"lastname" json:"lastName"`
	Login     string `db:"login" json:"login"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`

	// настраивается в профиле
	Gender   string `db:"gender" json:"gender"`
	Birthday string `db:"birthday" json:"birthday"`
	Location string `db:"location" json:"location"`
	Job      string `db:"job" json:"job"`
	Phone    string `db:"phone" json:"phone"`
	Links    string `db:"links" json:"links"`
	Avatar   string `db:"avatar" json:"avatar"`
	Layout   string `db:"layout" json:"layout"`
	Bio      string `db:"bio" json:"bio"`
}
