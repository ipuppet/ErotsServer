package logic

import (
	"ErotsServer/app/admin/dao"
	userDao "ErotsServer/app/user/dao"
	userLogic "ErotsServer/app/user/logic"
)

func NewUser(uid int) *dao.User {
	user := &dao.User{
		User: &userDao.User{},
	}
	user.Uid = uid

	return user
}

func GetUsers(page int, count int) []map[string]interface{} {
	start := (page - 1) * count
	end := count

	users, err := dao.GetUsers(start, end)
	if err != nil {
		return []map[string]interface{}{}
	}

	return users
}

func SearchUser(count int, kw string) []map[string]interface{} {
	users, err := dao.SearchUser(count, kw)
	if err != nil {
		return []map[string]interface{}{}
	}

	return users
}

func UpdateUserInfo(uid int, info map[string]interface{}) error {
	user := NewUser(uid)
	return user.UpdateInfo(info, userLogic.AdminEdit)
}

func GetUserRoles(uid int) ([]userDao.Role, error) {
	user := NewUser(uid)

	var roles []userDao.Role
	roles, err := user.GetRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func AddUserRole(uid int, roleId int) error {
	user := NewUser(uid)
	return user.AddRole(roleId)
}

func DeleteUserRole(uid int, roleId int) error {
	user := NewUser(uid)
	return user.DeleteRole(roleId)
}
