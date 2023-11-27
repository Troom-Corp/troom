package models

type User struct {
	UserId int    `db:"userid" json:"id"`
	Role   string `db:"role" json:"-"`

	// первый этап регистрации
	FirstName string `db:"firstname" json:"firstName"`
	LastName  string `db:"secondname" json:"lastName"`
	Login     string `db:"nick" json:"login"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`

	// второй этап регистрации
	Gender   string `db:"gender" json:"gender"`
	Birthday string `db:"dateofbirth" json:"birthday"`
	Location string `db:"location" json:"location"`
	Job      string `db:"job" json:"job"`

	// настраивается в профиле
	Phone  string `db:"phone" json:"phone"`   // in profile
	Links  string `db:"links" json:"links"`   // in profile
	Avatar string `db:"avatar" json:"avatar"` // in profile
	Bio    string `db:"bio" json:"bio"`       // in profile
}
