package logic

import "github.com/golang-jwt/jwt"

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
