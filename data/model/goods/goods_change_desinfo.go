package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type GoodsChangeDesInfo struct {
	Id         int    `json:"id"`
	ChangeId   int    `json:"change_id"`
	ChangeType string `json:"change_type"`
	// 产品编码
	CpCode string `json:"cp_code"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`
}

const GoodsChangeDesInfoTableName = "goods_change_des_info"

type GoodsChangeDesInfoList []GoodsChangeDesInfo
type GoodsChangeDesInfoDtoList []GoodsChangeDesInfoDto
type GoodsChangeDesInfoDto struct {
	Id         int    `json:"id"`
	ChangeId   int    `json:"change_id"`
	ChangeType string `json:"change_type"`
	// 产品编码
	CpCode string `json:"cp_code"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`
	// * ？、/
	Types   string  `json:"types"`
	ValuesL float64 `json:"values_l"`
}

func (r *GoodsChangeDesInfo) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDesInfoTableName).Create(r).Error
	return err
}
func (r *GoodsChangeDesInfo) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsChangeDesInfoTableName).Where("id = ?", r.Id).Delete(r).Error
}
func (r *GoodsChangeDesInfo) DeleteByType(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsChangeDesInfoTableName).Where("cp_code = ? and change_type = ?", r.CpCode, r.ChangeType).Delete(r).Error
}

func (r *GoodsChangeDesInfo) DeleteByChange_id(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsChangeDesInfoTableName).Where("change_id = ? ", r.ChangeId).Delete(r).Error
}

func (r *GoodsChangeDesInfo) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDesInfoTableName).Where("cp_code=?", cpType).Find(&r).Error
	return err
}

func (r *GoodsChangeDesInfo) GetByTypCode(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDesInfoTableName).Where("cp_code=? and change_type = ? ", r.CpCode, r.ChangeType).Find(&r).Error
	return err
}
func (r *GoodsChangeDesInfo) GetById(d *gorm.DB, id int) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDesInfoTableName).Where("id=?", id).Find(&r).Error
	return err
}
func (rList *GoodsChangeDesInfoDtoList) GetListByCpCode(cpCode string, d *gorm.DB) (err error) {
	if d == nil {
		d = db.BaseDB
	}

	sql := fmt.Sprintf("select b.id,a.change_id,a.change_type,a.cp_code,a.create_time,b.types,b.values_l"+
		" from goods_change_des_info a"+
		" inner join goods_change_desc b "+
		"on a.change_id = b.id where cp_code = %s", cpCode)

	err = d.Raw(sql).Find(&rList).Error
	return

}

func (r *GoodsChangeDesInfo) Save(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsChangeDesInfoTableName).Save(r).Error
	return err
}
