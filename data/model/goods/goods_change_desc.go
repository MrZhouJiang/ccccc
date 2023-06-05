package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type GoodsChangeDesc struct {
	Id    int    `json:"id"`
	CName string `json:"c_name"`
	// 1 损耗 2 换算
	ChangeType string `json:"change_type"`
	// * ？、/
	Types   string  `json:"types"`
	ValuesL float64 `json:"values_l"`

	CreateTime time.Time `json:"create_time"`
}

const GoodsChangeDescTableName = "goods_change_desc"

type GoodsChangeDescList []GoodsChangeDesc

func (r *GoodsChangeDesc) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDescTableName).Create(r).Error
	return err
}

func (r *GoodsChangeDesc) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDescTableName).Where("id=?", r.Id).Delete(r).Error
	return err
}

func (r *GoodsChangeDesc) Save(d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	r.CreateTime = time.Now()
	err := d.Table(GoodsChangeDescTableName).Save(r).Error
	return err
}

func (r *GoodsChangeDesc) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDescTableName).Where("cp_code=?", cpType).Find(&r).Error
	return err
}
func (r *GoodsChangeDesc) GetById(d *gorm.DB, id int) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDescTableName).Where("id=?", id).Find(&r).Error
	return err
}
func (r *GoodsChangeDesc) GetByName(d *gorm.DB, name string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDescTableName).Where("c_name=?", name).Find(&r).Error
	return err
}
func (rList *GoodsChangeDescList) GetListPage(changeType, name string, offset, limit int, d *gorm.DB) (total int, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(GoodsChangeDescTableName)
	if changeType != "" {
		//说明查询条件不为空
		query = query.Where("change_type = ? ", changeType)
	}
	if name != "" {
		//说明查询条件不为空
		query = query.Where("c_name like ? ", fmt.Sprintf("%%%s%%", name))
	}

	query.Count(&total)
	err = query.Order("id desc").Offset(offset).Limit(limit).Find(&rList).Error
	return

}
