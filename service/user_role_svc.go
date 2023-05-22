package service

import model_user "ccccc/data/model/user"
import model_role "ccccc/data/model/role"

//获取用户列表
func GetUserList() {
	list := model_user.UserInfoList{}
	err := list.GetList()
	if err != nil {

	}
}

//创建用户
func CreateUser(info model_user.UserInfo) {
	err := info.Create(nil)
	if err != nil {

	}

}

//删除用户
func DeleteUser(userId string) {
	info := model_user.UserInfo{
		UserId: userId,
	}
	err := info.Delete(nil)
	if err != nil {

	}
}

//根据用户ID 或者昵称获取用户基本信息
func GetUserById(userId, nick string) model_user.UserInfo {
	info := model_user.UserInfo{
		UserId: userId,
	}
	err := info.Get(userId, nick)
	if err != nil {
	}
	return info

}

type UserPermission struct {
	UserId   string
	Name     string
	RoleId   int
	RoleName string
}

//获取用户权限
func GetUserRole(userId, nick string) model_user.UserInfo {
	info := model_user.UserInfo{
		UserId: userId,
	}
	err := info.Get(userId, nick)
	if err != nil {
	}

	return info

}

func GetRoleList() {
	list := model_role.RoleList{}
	err := list.GetList()
	if err != nil {

	}
}

func GetPermissionList() {
	list := model_role.PermissionList{}
	err := list.GetList()
	if err != nil {

	}

}

// 获取权限列表
