package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const FenWeiTableName = "fen_wei"

type FenWei struct {
	FWBM       string    `json:"fwbm"`
	FWMC       string    `json:"fwmc"`
	FWXS       float64   `json:"fwxs"`
	FWJC       string    `json:"fwjc"`
	CreateUser string    `json:"create_user"`
	CreateTime time.Time `json:"create_time"`
}

type FenWeiList []FenWei

func (r *FenWei) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(FenWeiTableName).Create(r).Error
	return err
}

func (r *FenWei) DeleteByName(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(FenWeiTableName).Where("fwmc=? a", r.FWMC).Delete(r).Error
	return err
}

func (r *FenWei) GetByName(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(FenWeiTableName).Where("fwmc=? ", r.FWMC).Find(&r).Error
	return err
}

func (r *FenWei) GetByCode(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(FenWeiTableName).Where("fwbm=? ", r.FWBM).Find(&r).Error
	return err
}

func (r *FenWeiList) GetAll(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(FenWeiTableName).Where("1=1").Find(&r).Error
	return err
}
