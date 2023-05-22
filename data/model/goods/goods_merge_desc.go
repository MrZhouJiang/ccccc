package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

type GoodMergeDesc struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	CLName     string    `json:"cl_name"`
	Unit       string    `json:"unit"`
	Price      float64   `json:"price"`
	CreateTime time.Time `json:"create_time"`
}

const GoodMergeDescTableName = "goods_merge_desc"

type GoodMergeDescList []GoodMergeDesc

func (r *GoodMergeDesc) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Create(r).Error
	return err
}
func (r *GoodMergeDesc) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Where("id=?", r.Id).Delete(r).Error
	return err
}
func (r *GoodMergeDesc) Save(d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Save(r).Error
	return err
}

func (r *GoodMergeDesc) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Where("cp_code=?", cpType).Find(&r).Error
	return err
}
func (r *GoodMergeDesc) GetById(d *gorm.DB, id int) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Where("id=?", id).Find(&r).Error
	return err
}
func (r *GoodMergeDesc) GetByName(d *gorm.DB, name string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(GoodMergeDescTableName).Where("c_name=?", name).Find(&r).Error
	return err
}
func (rList *GoodMergeDescList) GetListPage(name string, offset, limit int, d *gorm.DB) (total int, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(GoodMergeDescTableName)
	if name != "" {
		//说明查询条件不为空
		query = query.Where("name = ? ", name)
	}

	query.Count(&total)
	err = query.Order("id desc").Offset(offset).Limit(limit).Find(&rList).Error
	return

}
