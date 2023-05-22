package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
)

type UnitDesc struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Types int    `json:"types"`
}

const UnitDescTableName = "unit_desc"

type UnitDescList []UnitDesc

func (r *UnitDesc) Create(d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(UnitDescTableName).Create(r).Error
	return err
}

func (r *UnitDesc) GetById(d *gorm.DB, id int) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(UnitDescTableName).Where("id=?", id).Find(&r).Error
	return err
}

func (r *UnitDescList) GetAll(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(UnitDescTableName).Where("1=1").Find(&r).Error
	return err
}
