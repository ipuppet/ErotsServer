package database

import (
	"ErotsServer/app/admin/internal"

	"github.com/ipuppet/gtools/database"
)

func NewUser(uid int) *internal.User {
	user := &internal.User{}
	user.Db = Db
	user.Uid = uid

	return user
}

func removeUserPrivateKey(user *map[string]interface{}) {
	delete(*user, "password")
}

func GetUsers(start int, end int) []map[string]interface{} {
	users, err := database.SQLQueryRetrieveMap(Db, `select * from ums.user limit ?,?`, start, end)
	if err != nil {
		return []map[string]interface{}{}
	}

	for index := range users {
		removeUserPrivateKey(&users[index])
	}

	return users
}

func SearchUser(count int, kw string) []map[string]interface{} {
	users, err := database.SQLQueryRetrieveMap(Db, "select * from ums.user where username like ? limit 0,?",
		"%"+kw+"%", count)
	if err != nil {
		return []map[string]interface{}{}
	}

	for index := range users {
		removeUserPrivateKey(&users[index])
	}

	return users
}
