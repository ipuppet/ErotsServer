package database

import (
	"ErotsServer/app/admin/structure"

	"github.com/ipuppet/gtools/database"
)

func GetRoles() ([]map[string]interface{}, error) {
	return database.SQLQueryRetrieveMap(Db, "select * from ums.rbac_role")
}

func UpdateRole(role structure.Role) error {
	_, err := database.MustExec(Db.Exec(`update ums.rbac_role
		set role_id=?, name=?, description=?
		where role_id=?`,
		*role.RoleId, role.Name, role.Description, *role.RoleId))

	return err
}

func AddRole(role structure.Role) error {
	_, err := database.MustExec(Db.Exec(`insert into ums.rbac_role
		(role_id, name, description)
		values
		(?, ?, ?)`,
		*role.RoleId, role.Name, role.Description))

	return err
}

func DeleteRole(roleId int) error {
	conn, err := Db.Begin()
	if err != nil {
		Logger.Println(err)
		return err
	}

	_, err = database.MustExec(conn.Exec("delete from ums.rbac_role where role_id=?", roleId))
	if err != nil {
		conn.Rollback()
		return err
	}

	_, err = database.MustExec(conn.Exec("delete from ums.rbac_role_permission where role_id=?", roleId))
	if err != nil {
		conn.Rollback()
		return err
	}

	conn.Commit()

	return nil
}

func GetPermissions() ([]map[string]interface{}, error) {
	return database.SQLQueryRetrieveMap(Db, "select * from ums.rbac_permission")
}

func UpdatePermission(permission structure.Permission) error {
	_, err := database.MustExec(Db.Exec(`update ums.rbac_permission
		set permission_id=?, module=?, description=?
		where permission_id=?`,
		*permission.PermissionId, permission.Name, permission.Description, *permission.PermissionId))

	return err
}

func AddPermission(permission structure.Permission) error {
	_, err := database.MustExec(Db.Exec(`insert into ums.rbac_permission
		(permission_id, module, description)
		values
		(?, ?, ?)`,
		*permission.PermissionId, permission.Name, permission.Description))

	return err
}

func DeletePermission(permissionId int) error {
	conn, err := Db.Begin()
	if err != nil {
		Logger.Println(err)
		return err
	}

	_, err = database.MustExec(conn.Exec("delete from ums.rbac_permission where permission_id=?", permissionId))
	if err != nil {
		conn.Rollback()
		return err
	}

	_, err = database.MustExec(conn.Exec("delete from ums.rbac_role_permission where permission_id=?", permissionId))
	if err != nil {
		conn.Rollback()
		return err
	}

	conn.Commit()

	return nil
}

func GetRolePermissions(roleId int) ([]map[string]interface{}, error) {
	return database.SQLQueryRetrieveMap(Db,
		`select a.permission_id,b.module,b.description
		from ums.rbac_role_permission a
		left join ums.rbac_permission b on a.role_id=?
		where a.permission_id=b.permission_id`,
		roleId)
}

func DeleteRolePermission(roleId int, permissionId int) error {
	_, err := database.MustExec(Db.Exec(`delete from ums.rbac_role_permission where role_id=? and permission_id=?`,
		roleId, permissionId))

	return err
}

func AddRolePermission(roleId int, permissionId int) error {
	_, err := database.MustExec(Db.Exec(`insert into ums.rbac_role_permission
		(role_id, permission_id)
		values
		(?, ?)`,
		roleId, permissionId))

	return err
}
