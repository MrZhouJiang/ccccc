package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

type GoodsType struct {
	//ID
	GoodsTypeId string `json:"goods_type_id"`
	//编码
	GoodsTypeCode string `json:"goods_type_code"`
	//名称
	GoodsTypeName string `json:"goods_type_name"`
	//创建者
	CreateUser string `json:"create_user"`
	//创建时间
	CreateTime time.Time `json:"create_time"`
	// 1 沙发成品 2 原材料
	Types int `json:"types"`
}

const GoodsTypeTableName = "goods_type"

type GoodsTypeList []GoodsType

func (r *GoodsType) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsTypeTableName).Create(r).Error
	return err
}

func (r *GoodsType) DeleteByName(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsTypeTableName).Where("goods_type_name=? ", r.GoodsTypeName).Delete(r).Error
	return err
}

func (r *GoodsType) GetByName(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsTypeTableName).Where("goods_type_name=? ", r.GoodsTypeName).Find(&r).Error
	return err
}

func (r *GoodsType) GetById(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsTypeTableName).Where("goods_type_id=? ", r.GoodsTypeId).Find(&r).Error
	return err
}

func (r *GoodsTypeList) GetAll(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsTypeTableName).Where("1=1").Find(&r).Error
	return err
}
