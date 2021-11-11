package model

import "fmt"

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l LoginData) IsValid() bool {
	return l.Username != ""
}

type RegisteredUser struct {
	Id           string
	Username     string
	PasswordHash string `pg:"password"`
}

func (r *RegisteredUser) String() string {
	return fmt.Sprintf("RegisteredUser<%s, %s, %s>", r.Id, r.Username, r.PasswordHash)
}
