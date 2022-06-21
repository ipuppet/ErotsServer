package pkg

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Uid            int
	Username       string
	Nickname       string
	Email          string
	Phone          interface{}
	Avatar         string
	Sex            int
	Password       string
	Lock           int
	RegisteredDate time.Time
	LastLoginDate  time.Time
}

type AccessTokenClaims struct {
	Uid   int
	Roles []map[string]interface{}
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	Uid int
	Ip  string
	jwt.StandardClaims
}
