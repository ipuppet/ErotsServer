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
)

func UpdateInfo(user *dao.User, userInfo map[string]interface{}) error {
	return user.UpdateInfo(userInfo, UserEdit)
}
