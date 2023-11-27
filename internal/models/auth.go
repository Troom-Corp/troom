package models

type IsCredentials struct {
	Login string `json:"isLogin"`
	Email string `json:"isEmail"`
}

//func (is IsCredentials) Validate() bool {
//	if is.Login == "" || is.Email == "" {
//		return false
//	}
//	return true
//}

type SignInCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpCredentials struct {
	FirstName  string `json:"firstName"`
	SecondName string `json:"lastName"`
	Login      string `json:"login"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Gender     string `json:"gender"`
	Birthday   string `json:"birthday"`
	Location   string `json:"location"`
	Job        string `json:"job"`
}

func (s SignUpCredentials) Validate() bool {
	if s.FirstName == "" || s.SecondName == "" || s.Login == "" || s.Email == "" || s.Password == "" || s.Gender == "" || s.Birthday == "" || s.Location == "" || s.Job == "" {
		return false
	}
	return true
}
