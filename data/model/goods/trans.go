package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

type Trans struct {
	Id        int    `json:"id"`
	TransId   string `json:"trans_id"`
	ShafaCode string `json:"shafa_code"`
	// 0 初始状态
	// 1 当前已经上线 且没有新的 草稿
	//  2 当前有草稿状态
	IsSubmit   int       `json:"is_submit"`
	CreateTime time.Time `json:"create_time"`
	OnlineTime time.Time `json:"online_time"`
	OnlineUser string    `json:"online_user"`
	CreateUser string    `json:"create_user"`
}

const TransTableName = "trans"

type TransList []Trans

func (r *Trans) Create(d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(TransTableName).Create(r).Error
	return err
}
func (r *Trans) Update(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(TransTableName).Save(r).Error
	return err
}

func (r *TransList) GetByShafaId(d *gorm.DB, id string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(TransTableName).Where("shafa_code=?", id).Order("create_time desc").Find(&r).Error
	return err
}

func (r *Trans) GetByShafaIdAndTrans(d *gorm.DB, id, trans string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(TransTableName).Where("shafa_code=? and trans_id = ?", id, trans).Order("create_time desc").Find(&r).Error
	return err
}

func (r *TransList) GetAll(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(TransTableName).Where("1=1").Find(&r).Error
	return err
}
