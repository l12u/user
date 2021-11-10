package model

type LoginData struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

func (l LoginData) IsValid() bool {
	return l.Username != ""
}
