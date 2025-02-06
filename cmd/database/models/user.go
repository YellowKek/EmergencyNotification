package models

type User struct {
	Id       int32  `json:"id"`
	Email    string `json:"email"`
	password string
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Rename   []string `json:"rename"` // список адресов почт и всего такого куда в случае чего отправлять
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}
