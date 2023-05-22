package db

import (
	"github.com/jinzhu/gorm"
)

type Test struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

const sdddddd = "ttt"

func (r *Test) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = Sqlserver_db
	}
	err := d.Table(sdddddd).Where("1=1").Find(&r).Error
	return err
}
