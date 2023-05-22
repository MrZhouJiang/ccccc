package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type ShaFaImportLog struct {
	Id         int       `json:"id"`
	SfName     string    `json:"sf_name"`
	SfCode     string    `json:"sf_code"`
	SDesc      string    `json:"s_desc"`
	ImportUser string    `json:"import_user"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	//沙发规格（可以解析分位名称）
	GG string `json:"gg"`
	//是否计算过
	IsSums string `json:"is_sums"`
}

const import_tabel_name = "sha_fa_import_log"

type ShaFaImportLogList []ShaFaImportLog

func (r *ShaFaImportLog) Create(d *gorm.DB) error {
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(import_tabel_name).Create(r).Error
	return err
}
func (r *ShaFaImportLog) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(import_tabel_name).Where("sf_code=?", r.SfCode).Delete(r).Error
	return err
}

func (r *ShaFaImportLog) GetByType(d *gorm.DB, sf_code string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(import_tabel_name).Where("sf_code=?", sf_code).Find(&r).Error
	return err
}

func (rList *ShaFaImportLogList) GetListPage(name, code string, offset, limit int, issums string, d *gorm.DB) (total int, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(import_tabel_name)
	if name != "" {
		//说明查询条件不为空
		query = query.Where("sf_name like ? ", fmt.Sprintf("%%%s%%", name))
	}

	if code != "" {
		//说明查询条件不为空
		query = query.Where("sf_code like ? ", fmt.Sprintf("%%%s%%", code))
	}
	if issums != "" {
		//说明查询条件不为空
		query = query.Where("is_sums = ? ", issums)
	}

	query.Count(&total)
	err = query.Order("create_time desc").Offset(offset).Limit(limit).Find(&rList).Error
	return

}

func (r *ShaFaImportLog) Get(name, code string) error {
	if name != "" {
		err := db.BaseDB.Table(import_tabel_name).Where("sf_name= ?", name).Find(&r).Error
		return err
	}
	if code != "" {
		err := db.BaseDB.Table(import_tabel_name).Where("sf_code= ?", code).Find(&r).Error
		return err
	}
	return nil
}

func (r *ShaFaImportLog) Update(d *gorm.DB) error {
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(import_tabel_name).Save(r).Error
	return err
}
