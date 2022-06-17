package database

import (
	"strings"

	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/database"
)

type User struct {
	Uid         int
	Roles       []map[string]interface{}
	Permissions map[int]interface{}
}

func (user *User) GetRoles() []map[string]interface{} {
	if len(user.Roles) == 0 {
		result, _ := database.SQLQueryRetrieveMap(Db,
			`select b.name,c.role_id
			from ums.rbac_role b,(select role_id from ums.rbac_user_role where uid=?) c
			where b.role_id=c.role_id`,
			user.Uid)
		user.Roles = result
	}

	return user.Roles
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

func (user *User) GetPermittedModules() []string {
	if len(user.Permissions) == 0 {
		user.Permissions = map[int]interface{}{}
		for _, role := range user.Roles {
			permission, _ := database.SQLQueryRetrieveMap(Db,
				`select a.role_id,a.permission_id,b.module,b.description
				from ums.rbac_role_permission a
				left join ums.rbac_permission b on a.role_id=?
				where a.permission_id=b.permission_id`,
				role["role_id"])

			if len(permission) > 0 {
				user.Permissions[int(role["role_id"].(float64))] = permission
			}
		}
	}

	modules := []string{}
	tmp := make(map[int64]byte, 10) // 去重用到的临时 map

	for _, permissions := range user.Permissions {
		for _, permission := range permissions.([]map[string]interface{}) {
			l := len(tmp)
			tmp[permission["permission_id"].(int64)] = 0
			if l != len(tmp) {
				modules = append(modules, permission["module"].(string))
			}
		}
	}

	return modules
}

func (user *User) UpdateInfo(info map[string]interface{}) error {
	userStructure := map[string]interface{}{}
	config.GetConfig("userStructure.json", &userStructure)
	canEdit := userStructure["canEdit"].([]interface{})
	adminEdit := userStructure["adminEdit"].([]interface{})
	canEdit = append(canEdit, adminEdit...)

	columnList := []string{}
	valueList := []interface{}{}

	for _, v := range canEdit {
		column := v.(string)
		value, ok := info[column]
		if ok {
			columnList = append(columnList, "`"+column+"`=?")
			valueList = append(valueList, value)
		}
	}
	valueList = append(valueList, user.Uid)

	sql := "update ums.user set " + strings.Join(columnList, ",") + " where `uid`=?"
	_, err := database.MustExec(Db.Exec(sql, valueList...))
	return err
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
