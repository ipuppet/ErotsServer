package logic

import (
	"ErotsServer/app/user/dao"
)

var (
	UserEdit = []string{
		"avatar",
		"nickname",
		"sex",
	}
	AdminEdit = append(UserEdit, []string{
		"lock",
	}...)
)

func UpdateInfo(user *dao.User, userInfo map[string]interface{}) error {
	if isPermitted, _ := user.IsPermitted("rbacManager"); isPermitted {
		return user.UpdateInfo(userInfo, AdminEdit)
	}

	return user.UpdateInfo(userInfo, UserEdit)
}
