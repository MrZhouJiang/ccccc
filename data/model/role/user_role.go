package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	role_tableName = "user_role"
)

//用户角色表
type UserRole struct {
	Id       int       `json:"id"`
	UserId   string    `json:"user_id"`
	RoleId   int       `json:"role_id"`
	RoleName string    `json:"role_name"`
	AddTime  time.Time `json:"add_time"`
	AddUser  string    `json:"add_user"`
	DescInfo string    `json:"desc_info"`
}
type UserPermission struct {
	RoleId       int    `json:"role_id"`
	PermissionId int    `json:"permission_id"`
	Type         int    `json:"type"`
	Url          string `json:"url"`
	Name         string `json:"name"`
}

func (r *UserRole) Create(d *gorm.DB) error {
	r.AddTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(role_tableName).Create(r).Error
	return err
}

type UserRoleList []UserRole
type UserPermissionList []UserPermission

func (list *UserRole) GetUserRole(userId string, d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}

	err := d.Table(role_tableName).Where("user_id = ?", userId).Find(&list).Error
	return err

}

//查询用户拥有的权限 （需要缓存）
func (role *UserRole) GetUserPermission() (list UserPermissionList, err error) {
	sql := fmt.Sprintf("select  role.role_id,permission.id as permission_id, permission.type,permission.url,permission.name  from user_role role "+
		" inner join  role_permission_relation as re on role.role_id =re.role_id "+
		" inner join  permission on permission.id = re.permission_id"+
		"  where user_id ='%s' ", role.UserId)
	query := db.BaseDB.Table(role_tableName)
	err = query.Raw(sql).Find(&list).Error
	return

}
