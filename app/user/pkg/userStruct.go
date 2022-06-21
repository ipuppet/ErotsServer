package pkg

import (
	"database/sql"
	"strings"

	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/database"
)

type User struct {
	Db          *sql.DB
	Uid         int
	Roles       []map[string]interface{}
	Permissions map[int]interface{}
}

func (user *User) GetRoles() []map[string]interface{} {
	if len(user.Roles) == 0 {
		result, _ := database.SQLQueryRetrieveMap(user.Db,
			`select b.name,c.role_id
			from ums.rbac_role b,(select role_id from ums.rbac_user_role where uid=?) c
			where b.role_id=c.role_id`,
			user.Uid)
		user.Roles = result
	}

	return user.Roles
}

func (user *User) GetPermittedModules() []string {
	if len(user.Permissions) == 0 {
		user.Permissions = map[int]interface{}{}
		for _, role := range user.Roles {
			permission, _ := database.SQLQueryRetrieveMap(user.Db,
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
	_, err := database.MustExec(user.Db.Exec(sql, valueList...))
	return err
}
