package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

type GongYi struct {
	Id      int    `json:"id"`
	ShafaId string `json:"shafa_id"`
	//裁工
	Types      string `json:"types"`
	FenWeiName string `json:"fen_wei_name"`
	CLName     string `json:"cl_name"`
	Size       string `json:"size"`
	Nums       string `json:"nums"`
	Unit       string `json:"unit"`
	Descs      string `json:"descs"`
	TotalPrice string `json:"total_price"`
	CpCode     string `json:"cp_code"`
	//大类 工艺名称
	GongYiName string    `json:"gong_yi_name"`
	CreateTime time.Time `json:"create_time"`
	JiJiaNum   string    `json:"ji_jia_num"`
	//使用部位
	GoodsPoint string `json:"goods_point"`
}

const gongyiTableName = "sofa_gong_yi"

type GongYiList []GongYi

func (r *GongYi) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyiTableName).Create(r).Error
	return err
}
func (r *GongYi) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyiTableName).Where("shafa_id=?", r.ShafaId).Delete(r).Error
	return err
}

func (r *GongYi) DeleteByTypes(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyiTableName).Where("shafa_id=? and types = ?", r.ShafaId, r.Types).Delete(r).Error
	return err
}

func (r *GongYi) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyiTableName).Where("cp_type=?", cpType).Find(&r).Error
	return err
}

func (rList *GongYiList) GetListPage(shafaId, types string, d *gorm.DB) (err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(gongyiTableName)
	if shafaId != "" {
		//说明查询条件不为空
		query = query.Where("shafa_id  =? ", shafaId)
	}

	if types != "" {
		//说明查询条件不为空
		query = query.Where("types  =? ", types)
	}

	err = query.Order("gong_yi_name").Find(&rList).Error
	return

}

func (r *GongYiList) GetBySoFaCode(d *gorm.DB, sofaCode string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyiTableName).Where("shafa_id=?", sofaCode).Find(&r).Error
	return err
}
