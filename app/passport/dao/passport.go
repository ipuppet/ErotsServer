package dao

import (
	"bytes"
	"regexp"
	"time"

	"github.com/ipuppet/gtools/database"
	"github.com/ipuppet/gtools/regex"
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

func GetUserInfo(account string) (User, error) {
	var column string = "username"
	accountByte := []byte(account)
	if matched, _ := regexp.Match(regex.Email, accountByte); matched {
		column = "email"
	} else if matched, _ := regexp.Match(regex.Phone_cn, accountByte); matched {
		column = "phone"
	}

	user := User{}

	row := Db.QueryRow(getUserInfoSQL(column), account)
	ScanUserRow(row, &user)

	return user, nil
}

func GetUserInfoByUid(uid int, user *User) error {
	row := Db.QueryRow(getUserInfoSQL("uid"), uid)
	return ScanUserRow(row, user)
}

func GetUserRolesByUid(uid int) []map[string]interface{} {
	result, _ := database.SQLQueryRetrieveMapNoCache(Db,
		`select b.name,c.role_id
		from ums.rbac_role b,(select role_id from ums.rbac_user_role where uid=?) c
		where b.role_id=c.role_id`,
		uid)

	return result
}

func LogLoginInfo(user User) {
	Db.Exec(`update ums.user set last_login_date=FROM_UNIXTIME(?) where uid=?`, time.Now().Unix(), user.Uid)
}

func Register(nickname string, username string, email string, password []byte) error {
	conn, err := Db.Begin()
	if err != nil {
		return err
	}

	result, err := database.MustExec(conn.Exec(
		`insert into ums.user
		(username, nickname, email, password)
		values
		(?, ?, ?, ?)`,
		username, nickname, email, password,
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
