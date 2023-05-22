package sysnc

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
)

const gongyiTableName = "BASE_CP"

type SqlCpBase struct {
	CPFLID   int     `gorm:"column:CPFLID" json:"CPFLID"`
	CPMC     string  `gorm:"column:CPMC" json:"CPMC"`
	JLDWID_Z int     `gorm:"column:JLDWID_Z" json:"JLDWID_Z"`
	JLDWID_F int     `gorm:"column:JLDWID_F" json:"JLDWID_F"`
	CPBM     string  `gorm:"column:CPBM" json:"CPBM"`
	GG       string  `gorm:"column:GG" json:"GG"`
	CPJC     string  `gorm:"column:CPJC" json:"CPJC"`
	XS_ZJL   float64 `gorm:"column:XS_ZJL" json:"XS_ZJL"`
	XS_FJL   float64 `gorm:"column:XS_FJL" json:"XS_FJL"`
	ZHL      string  `gorm:"column:ZHL" json:"ZHL"`
	CBJ      float64 `gorm:"column:CBJ" json:"CBJ"`
	//通用名称
	TYMC string `gorm:"column:TYMC" json:"TYMC"`
	SFTY int    `gorm:"column:SFTY" json:"SFTY"`

	CJSJ string `gorm:"column:CJSJ" json:"CJSJ"`
	XGSJ string `gorm:"column:XGSJ" json:"XGSJ"`
}

type CpList []SqlCpBase

func (rList *CpList) GetListPage(offset, size int, d *gorm.DB) (err error) {
	if d == nil {
		d = db.SqlDB
	}
	query := d.Table(gongyiTableName).Select("CPFLID,CPMC,JLDWID_Z,JLDWID_F,CPBM,GG,CPJC,XS_ZJL,XS_FJL,ZHL,CBJ,TYMC,SFTY,CJSJ,XGSJ")

	err = query.Offset(offset).Limit(size).Order("CJSJ asc").Find(&rList).Error
	return

}

func (rList *SqlCpBase) GetListOne(d *gorm.DB) (err error) {
	if d == nil {
		d = db.SqlDB
	}
	query := d.Table(gongyiTableName)

	err = query.Where("1=1").Order("CJSJ desc").First(&rList).Error
	return

}
