package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type GoodsMergeDesInfo struct {
	Id      int `json:"id"`
	MergeId int `json:"merge_id"`
	// 产品编码
	CpCode string `json:"cp_code"`

	// 创建时间
	CreateTime time.Time `json:"create_time"`
}

const GoodsMergeDesInfoTableName = "goods_merge_des_info"

type GoodsMergeDesInfoList []GoodsMergeDesInfo
type GoodsMergeDesInfoDtoList []GoodsMergeDesInfoDto
type GoodsMergeDesInfoDto struct {
	Id      int `json:"id"`
	MergeId int `json:"merge_id"`
	// 产品编码
	CpCode     string    `json:"cp_code"`
	Name       string    `json:"name"`
	CLName     string    `json:"cl_name"`
	Unit       string    `json:"unit"`
	Price      float64   `json:"price"`
	CreateTime time.Time `json:"create_time"`
}

func (r *GoodsMergeDesInfo) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsMergeDesInfoTableName).Create(r).Error
	return err
}
func (r *GoodsMergeDesInfo) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsMergeDesInfoTableName).Where("id = ?", r.Id).Delete(r).Error
}
func (r *GoodsMergeDesInfo) DeleteByType(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsMergeDesInfoTableName).Where("cp_code = ? ", r.CpCode).Delete(r).Error
}

func (r *GoodsMergeDesInfo) DeleteByMerGetId(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	return d.Table(GoodsMergeDesInfoTableName).Where("merge_id = ? ", r.MergeId).Delete(r).Error
}

func (r *GoodsMergeDesInfo) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsMergeDesInfoTableName).Where("cp_code=?", cpType).Find(&r).Error
	return err
}

func (r *GoodsMergeDesInfo) GetByTypCode(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsMergeDesInfoTableName).Where("cp_code=?", r.CpCode).Find(&r).Error
	return err
}
func (r *GoodsMergeDesInfo) GetById(d *gorm.DB, id int) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodsMergeDesInfoTableName).Where("id=?", id).Find(&r).Error
	return err
}
func (rList *GoodsMergeDesInfoDtoList) GetListByCpCode(cpCode string, d *gorm.DB) (err error) {
	if d == nil {
		d = db.BaseDB
	}

	sql := fmt.Sprintf("select b.id,a.merge_id,a.cp_code,a.create_time,b.unit,b.price,b.cl_name,b.name"+
		" from goods_merge_des_info a"+
		" inner join goods_merge_desc b "+
		"on a.merge_id = b.id where cp_code = %s", cpCode)

	err = d.Raw(sql).Find(&rList).Error
	return

}

func (r *GoodsMergeDesInfo) Save(d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	r.CreateTime = time.Now()
	err := d.Table(GoodsMergeDesInfoTableName).Save(r).Error
	return err
}
