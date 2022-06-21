package database

import (
	"bytes"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"ErotsServer/app/passport/pkg"

	"github.com/golang-jwt/jwt"
	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/database"
	"github.com/ipuppet/gtools/regex"
	"github.com/ipuppet/gtools/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	Db         *sql.DB
	UserColumn string
)

func init() {
	Db = database.ConnectToMySQL("ums")

	columns, _ := utils.GetStructFieldNameToSnake(pkg.User{})
	UserColumn = "`" + strings.Join(columns, "`,`") + "`"
}

const (
	JwtIssuer        = "ErotsServer"
	accessTokenLife  = 24 * time.Hour
	refreshTokenLife = 7 * 24 * time.Hour
)

func getUserInfoSQL(column string) string {
	var sqlBuffer bytes.Buffer
	sqlBuffer.WriteString("select ")
	sqlBuffer.WriteString(UserColumn)
	sqlBuffer.WriteString(" from ums.user where ")
	sqlBuffer.WriteString(column)
	sqlBuffer.WriteString("=? limit 1")

	return sqlBuffer.String()
}

func getUserInfo(account string) (pkg.User, error) {
	var column string = "username"
	accountByte := []byte(account)
	if matched, _ := regexp.Match(regex.Email, accountByte); matched {
		column = "email"
	} else if matched, _ := regexp.Match(regex.Phone_cn, accountByte); matched {
		column = "phone"
	}

	user := pkg.User{}

	database.SQLQueryRetrieveStruct(Db, &user, getUserInfoSQL(column), account)

	return user, nil
}

func GetUserInfoByUid(uid int) pkg.User {
	user := pkg.User{}

	database.SQLQueryRetrieveStruct(Db, &user, getUserInfoSQL("uid"), uid)

	return user
}

func getUserRolesByUid(uid int) []map[string]interface{} {
	result, _ := database.SQLQueryRetrieveMapNoCache(Db,
		`select b.name,c.role_id
		from ums.rbac_role b,(select role_id from ums.rbac_user_role where uid=?) c
		where b.role_id=c.role_id`,
		uid)

	return result
}

func generateAccessToken(user pkg.User, roles []map[string]interface{}) (string, error) {
	expireTime := time.Now().Add(accessTokenLife)
	jwtClaims := pkg.AccessTokenClaims{
		Uid:   user.Uid,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    JwtIssuer,         // 签名颁发者
			Subject:   "UserAccessToken", // 签名主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(utils.MD5(user.Password)))
}

func generateRefreshToken(user pkg.User, ip string) (string, error) {
	expireTime := time.Now().Add(refreshTokenLife)
	jwtClaims := pkg.RefreshTokenClaims{
		Uid: user.Uid,
		Ip:  ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    JwtIssuer,          // 签名颁发者
			Subject:   "UserRefreshToken", // 签名主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString([]byte(utils.MD5(user.Password)))
}

func getUserMapForToken(user pkg.User, roles []map[string]interface{}) map[string]interface{} {
	userMap := utils.StructToMapWithLowerKey(user)

	userMap["roles"] = roles

	userStructure := map[string]interface{}{}
	config.GetConfig("userStructure.json", &userStructure)

	for _, column := range userStructure["cannotUseForSignUp"].([]interface{}) {
		delete(userMap, column.(string))
	}

	return userMap
}

func ByPassword(account string, password string, ip string) (map[string]interface{}, error) {
	user, err := getUserInfo(account)
	if err != nil {
		return nil, err
	}

	// password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("wrong password")
	}
	// lock
	if user.Lock != 0 {
		return nil, errors.New("account has been locked")
	}

	// user roles
	roles := getUserRolesByUid(user.Uid)

	accessToken, err := generateAccessToken(user, roles)
	if err != nil {
		return nil, errors.New("access token generation failed")
	}
	refreshToken, err := generateRefreshToken(user, ip)
	if err != nil {
		return nil, errors.New("refresh token generation failed")
	}

	// log last_login_date
	Db.Exec(`update ums.user set last_login_date=FROM_UNIXTIME(?) where uid=?`, time.Now().Unix(), user.Uid)

	return map[string]interface{}{
		"public":        getUserMapForToken(user, roles),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"exp":           time.Now().Add(accessTokenLife).Unix(), // 参考过期时间
	}, nil
}

func ParseToken(accessTokenString string, refreshTokenString string, ip string) (map[string]interface{}, error) {
	var user pkg.User
	var tokenKey string
	accessToken, err := jwt.ParseWithClaims(accessTokenString, &pkg.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*pkg.AccessTokenClaims)
		user = GetUserInfoByUid(claims.Uid)
		tokenKey = utils.MD5(user.Password)
		return []byte(tokenKey), nil
	})

	// accessToken 验证失败
	if err != nil {
		// refresh token
		refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &pkg.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			claims := token.Claims.(*pkg.RefreshTokenClaims)
			user = GetUserInfoByUid(claims.Uid)
			tokenKey = utils.MD5(user.Password)
			return []byte(tokenKey), nil
		})
		if err != nil {
			return nil, errors.New("refresh token parse failed: " + err.Error())
		}

		if claims, ok := refreshToken.Claims.(*pkg.RefreshTokenClaims); ok && refreshToken.Valid {
			if claims.Ip != ip {
				return nil, errors.New("refresh token verification failed: different ip addresse")
			}

			roles := getUserRolesByUid(user.Uid)

			// new access token
			accessToken, err := generateAccessToken(user, roles)
			if err != nil {
				return nil, errors.New("access token generation failed")
			}

			// new refresh token
			refreshToken, err := generateRefreshToken(user, ip)
			if err != nil {
				return nil, errors.New("refresh token generation failed")
			}

			return map[string]interface{}{
				"public":        getUserMapForToken(user, roles),
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			}, nil
		}
	} else {
		if claims, ok := accessToken.Claims.(*pkg.AccessTokenClaims); ok && accessToken.Valid {
			// 使用令牌内的信息，减少查库次数
			return map[string]interface{}{
				"public": getUserMapForToken(user, claims.Roles),
			}, nil
		}
	}

	return nil, err
}

func Register(nickname string, email string, password string) error {
	conn, err := Db.Begin()
	if err != nil {
		return err
	}

	username := utils.MD5(email)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	result, err := database.MustExec(conn.Exec(
		`insert into ums.user
		(username, nickname, email, password)
		values
		(?, ?, ?, ?)`,
		username, nickname, email, hashedPassword,
	))
	if err != nil {
		conn.Rollback()
		return err
	}

	uid, err := result.LastInsertId()
	if err != nil {
		conn.Rollback()
		return err
	}

	_, err = database.MustExec(conn.Exec(
		`insert into ums.rbac_user_role (uid, role_id) values (?, ?)`,
		uid, 1,
	))
	if err != nil {
		conn.Rollback()
		return err
	}

	return conn.Commit()
}
