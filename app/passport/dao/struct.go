package dao

import (
	"database/sql"
	"time"
)

type User struct {
	Uid            int
	Username       string
	Nickname       string
	Email          string
	Phone          string
	Avatar         string
	Sex            int
	Password       string
	Lock           int
	RegisteredDate time.Time
	LastLoginDate  time.Time
}

func ScanUserRow(row *sql.Row, user *User) error {
	return row.Scan(
		&user.Uid,
		&user.Username,
		&user.Nickname,
		&user.Email,
		&user.Phone,
		&user.Avatar,
		&user.Sex,
		&user.Password,
		&user.Lock,
		&user.RegisteredDate,
		&user.LastLoginDate,
	)
}
