package dao

import (
	"strings"

	"github.com/ipuppet/gtools/database"
)

type User struct {
	Uid              int
	roles            []Role
	permissions      []Permission
	permittedModules []string
}

func (user *User) GetRoles() ([]Role, error) {
	if len(user.roles) == 0 {
		rows, err := Db.Query(
			`select b.name,b.description,c.role_id
			from ums.rbac_role b,(select role_id from ums.rbac_user_role where uid=?) c
			where b.role_id=c.role_id`,
			user.Uid)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var role Role
			err := rows.Scan(&role.Name, &role.Description, &role.RoleId)
			if err != nil {
				Logger.Println(err)
				return nil, err
			}
			user.roles = append(user.roles, role)
		}

		if err = rows.Err(); err != nil {
			Logger.Println(err)
			return nil, err
		}
	}

	return user.roles, nil
}

func (user *User) GetPermissions() ([]Permission, error) {
	if len(user.permissions) == 0 {
		if len(user.roles) == 0 {
			_, err := user.GetRoles()
			if err != nil {
				return nil, err
			}
		}
		for _, role := range user.roles {
			rows, err := Db.Query(
				`select a.role_id,a.permission_id,b.module,b.description
				from ums.rbac_role_permission a
				left join ums.rbac_permission b on a.role_id=?
				where a.permission_id=b.permission_id`,
				role.RoleId)
			if err != nil {
				return nil, err
			}

			defer rows.Close()

			for rows.Next() {
				var permission Permission
				err := rows.Scan(&permission.RoleId, &permission.PermissionId, &permission.Module, &permission.Description)
				if err != nil {
					Logger.Println(err)
					return nil, err
				}
				user.permissions = append(user.permissions, permission)
			}

			if err = rows.Err(); err != nil {
				Logger.Println(err)
				return nil, err
			}
		}
	}

	return user.permissions, nil
}

func (user *User) GetPermittedModules() ([]string, error) {
	if len(user.permittedModules) == 0 {
		if len(user.permissions) == 0 {
			_, err := user.GetPermissions()
			if err != nil {
				return nil, err
			}
		}

		tmp := make(map[int]byte, 10) // 去重用到的临时 map

		for _, permission := range user.permissions {
			l := len(tmp)
			tmp[permission.PermissionId] = 0
			if l != len(tmp) {
				user.permittedModules = append(user.permittedModules, permission.Module)
			}
		}
	}

	return user.permittedModules, nil
}

func (user *User) IsPermitted(moduleName string) (bool, error) {
	permittedModules, err := user.GetPermittedModules()
	if err != nil {
		return false, err
	}

	for _, permittedModule := range permittedModules {
		if moduleName == permittedModule || permittedModule == "ALL" {
			return true, nil
		}
	}

	return false, nil
}

func (user *User) UpdateInfo(info map[string]interface{}, canEdit []string) error {
	columnList := []string{}
	valueList := []interface{}{}

	for _, column := range canEdit {
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
