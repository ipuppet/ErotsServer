package database

import (
	"database/sql"
	"log"

	passportDatabase "ErotsServer/app/passport/database"

	"github.com/golang-jwt/jwt"
	"github.com/ipuppet/gtools/database"
	"github.com/ipuppet/gtools/utils"
)

var (
	Db     *sql.DB
	Logger *log.Logger
)

func init() {
	Db = database.ConnectToMySQL("")
	Logger = utils.Logger("admin")
}

func ParseAccessToken(tokenString string) (*passportDatabase.AccessTokenClaims, error) {
	var user passportDatabase.User
	var tokenKey string
	accessToken, err := jwt.ParseWithClaims(tokenString, &passportDatabase.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*passportDatabase.AccessTokenClaims)
		user = passportDatabase.GetUserInfoByUid(claims.Uid)
		tokenKey = utils.MD5(user.Password)
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := accessToken.Claims.(*passportDatabase.AccessTokenClaims); ok && accessToken.Valid {
		// 使用令牌内的信息，减少查库次数
		return claims, nil
	}

	return nil, err
}
