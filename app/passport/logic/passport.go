package logic

import (
	"errors"
	"time"

	"ErotsServer/app/passport/dao"

	"github.com/golang-jwt/jwt"
	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	JwtIssuer        = "ErotsServer"
	accessTokenLife  = 24 * time.Hour
	refreshTokenLife = 7 * 24 * time.Hour
)

func generateAccessToken(user dao.User, roles []map[string]interface{}) (string, error) {
	expireTime := time.Now().Add(accessTokenLife)
	jwtClaims := AccessTokenClaims{
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

func generateRefreshToken(user dao.User, ip string) (string, error) {
	expireTime := time.Now().Add(refreshTokenLife)
	jwtClaims := RefreshTokenClaims{
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

func getUserMapForToken(user dao.User, roles []map[string]interface{}) map[string]interface{} {
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
	user, err := dao.GetUserInfo(account)
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
	roles := dao.GetUserRolesByUid(user.Uid)

	accessToken, err := generateAccessToken(user, roles)
	if err != nil {
		return nil, errors.New("access token generation failed")
	}
	refreshToken, err := generateRefreshToken(user, ip)
	if err != nil {
		return nil, errors.New("refresh token generation failed")
	}

	// log last_login_date
	go dao.LogLoginInfo(user)

	return map[string]interface{}{
		"public":        getUserMapForToken(user, roles),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"exp":           time.Now().Add(accessTokenLife).Unix(), // 参考过期时间
	}, nil
}

func ParseRefreshToken(tokenString string, user *dao.User) (*RefreshTokenClaims, error) {
	var tokenKey string
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*RefreshTokenClaims)
		if user == nil || user.Password == "" {
			dao.GetUserInfoByUid(claims.Uid, user)
		}
		tokenKey = utils.MD5(user.Password)
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		// 使用令牌内的信息，减少查库次数
		return claims, nil
	}

	return nil, err
}

func ParseAccessToken(tokenString string, user *dao.User) (*AccessTokenClaims, error) {
	var tokenKey string
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := token.Claims.(*AccessTokenClaims)
		if user == nil || user.Password == "" {
			dao.GetUserInfoByUid(claims.Uid, user)
		}
		tokenKey = utils.MD5(user.Password)
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		// 使用令牌内的信息，减少查库次数
		return claims, nil
	}

	return nil, err
}

func ParseToken(accessTokenString string, refreshTokenString string, ip string) (map[string]interface{}, error) {
	var user dao.User
	accessTokenClaims, err := ParseAccessToken(accessTokenString, &user)

	// accessToken 验证失败
	if err != nil {
		// refresh token
		refreshTokenClaims, err := ParseRefreshToken(refreshTokenString, &user)
		if err != nil {
			return nil, errors.New("refresh token parse failed: " + err.Error())
		}

		if refreshTokenClaims.Ip != ip {
			return nil, errors.New("refresh token verification failed: different ip addresse")
		}

		roles := dao.GetUserRolesByUid(user.Uid)

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

	// 使用令牌内的信息，减少查库次数
	return map[string]interface{}{
		"public": getUserMapForToken(user, accessTokenClaims.Roles),
	}, nil
}

func Register(nickname string, email string, password string) error {

	username := utils.MD5(email)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := dao.Register(nickname, username, email, hashedPassword)

	return err
}
