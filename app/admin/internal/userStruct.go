package internal

import (
	userPkg "ErotsServer/app/user/pkg"

	"github.com/ipuppet/gtools/database"
)

type User struct {
	*userPkg.User
}

func (user *User) AddRole(roleId int) error {
	_, err := database.MustExec(user.Db.Exec(`insert into ums.rbac_user_role (uid, role_id) values (?, ?)`,
		user.Uid, roleId))

	return err
}

func (user *User) DeleteRole(roleId int) error {
	_, err := database.MustExec(user.Db.Exec(`delete from ums.rbac_user_role where uid=? and role_id=?`,
		user.Uid, roleId))

	return err
}
