package models

type User struct {
	Username string `json:"username"`
	password []byte `json:"-"`
}

// SetPassword sets the user's password.
func (u *User) SetPassword(pass []byte) {
	u.password = pass
}

// GetPassword returns the hashed password.
func (u *User) GetPassword() []byte {
	return u.password
}

type UserSessions struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (UserSessions) TableName() string {
	return "usersessions"
}