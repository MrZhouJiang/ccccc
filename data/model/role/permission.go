package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	permission_tableName = "permission"
)

// 权限表
type Permission struct {
	Id      int       `json:"id"`
	Type    int       `json:"type"`
	Url     string    `json:"url"`
	Name    string    `json:"name"`
	AddTime time.Time `json:"add_time"`
}
type PermissionList []Permission

func (r *Permission) Create(d *gorm.DB) error {
	r.AddTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(permission_tableName).Create(r).Error
	return err
}

func (list *PermissionList) GetList() error {
	//剔除管理员账号

	return nil
}
