package models

type IsCredentials struct {
	Login string `json:"isLogin"`
	Email string `json:"isEmail"`
}

type SignInCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpCredentials struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (s SignUpCredentials) Validate() bool {
	if s.FirstName == "" || s.LastName == "" || s.Login == "" || s.Email == "" || s.Password == "" {
		return false
	}
	return true
}
