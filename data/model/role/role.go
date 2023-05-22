package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	role_tablename = "role"
)

//角色表
type Role struct {
	Id int `json:"id"`
	//角色名称
	RoleName string    `json:"role_name"`
	DescInfo string    `json:"desc_info"`
	AddTime  time.Time `json:"add_time"`
}

type RoleList []Role

func (r *Role) Create(d *gorm.DB) error {
	r.AddTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(role_tablename).Create(r).Error
	return err
}

func (list *RoleList) GetList() error {
	//剔除管理员账号

	return nil
}
