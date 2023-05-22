package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	role_relation_tableName = "role_permission_relation"
)

//角色权限表
type RolePermissionRelation struct {
	Id           int       `json:"id"`
	RoleId       int       `json:"role_id"`
	PermissionId int       `json:"permission_id"`
	AddTime      time.Time `json:"add_time"`
	AddUser      string    `json:"add_user"`
	DescInfo     string    `json:"desc_info"`
}

func (r *RolePermissionRelation) Create(d *gorm.DB) error {
	r.AddTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(role_relation_tableName).Create(r).Error
	return err
}
