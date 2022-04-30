package db

type User struct {
	Id           string `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Email        string `json:"email" db:"email"`
	Phone        string `json:"phone" db:"phone"`
}

type UserInRoles struct {
	UserId string `json:"user_id" db:"user_id"`
	RoleId string `json:"role_id" db:"role_id"`
}

type Role struct {
	Id     string `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Code   string `json:"code" db:"code"`
	Source string `json:"source" db:"source"`
}

type Session struct {
	UserId    string `json:"user_id" db:"user_id"`
	StartIn   int64  `json:"start_in" db:"start_in"`
	ExpiresIn int64  `json:"expires_in" db:"expired_in"`
	Token     string `json:"token" db:"token"`
}