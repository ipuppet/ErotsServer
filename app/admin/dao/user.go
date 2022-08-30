package dao

import (
	userDao "ErotsServer/app/user/dao"

	"github.com/ipuppet/gtools/database"
)

var (
	UserPrivateKey = []string{
		"password",
	}
)

type User struct {
	*userDao.User
}

func (user *User) AddRole(roleId int) error {
	_, err := database.MustExec(Db.Exec(`insert into ums.rbac_user_role (uid, role_id) values (?, ?)`,
		user.Uid, roleId))

	return err
}

func (user *User) DeleteRole(roleId int) error {
	_, err := database.MustExec(Db.Exec(`delete from ums.rbac_user_role where uid=? and role_id=?`,
		user.Uid, roleId))

	return err
}

func removeUserPrivateKey(user *map[string]interface{}) {
	for _, key := range UserPrivateKey {
		delete(*user, key)
	}
}

func GetUsers(start int, end int) ([]map[string]interface{}, error) {
	users, err := database.SQLQueryRetrieveMapNoCache(Db, `select * from ums.user limit ?,?`, start, end)
	if err != nil {
		return nil, err
	}

	for index := range users {
		removeUserPrivateKey(&users[index])
	}

	return users, nil
}

func SearchUser(count int, kw string) ([]map[string]interface{}, error) {
	users, err := database.SQLQueryRetrieveMapNoCache(Db, "select * from ums.user where username like ? limit 0,?",
		"%"+kw+"%", count)
	if err != nil {
		return nil, err
	}

	for index := range users {
		removeUserPrivateKey(&users[index])
	}

	return users, nil
}
