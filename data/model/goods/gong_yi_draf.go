package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

type GongYiDraf struct {
	Id      int    `json:"id"`
	ShafaId string `json:"shafa_id"`
	//裁工
	Types      string `json:"types"`
	FenWeiName string `json:"fen_wei_name"`
	CLName     string `json:"cl_name"`
	//填写的规格
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
	JiJiaUnit  string    `json:"ji_jia_unit"`
	//使用部位
	GoodsPoint   string  `json:"goods_point"`
	ShunHaoPrice float64 `json:"shun_hao_price"`
	ImportUser   string  `json:"import_user"`
	TransId      string  `json:"trans_id"`
	//固定的规格 默认是 基础表的 手动填写的话 就是手动填写的。
	OwnerSize string `json:"owner_size"`
}

const gongyidrafTableName = "sofa_draf_gong_yi"

type GongYiDrafList []GongYiDraf

func (r *GongYiDraf) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Create(r).Error
	return err
}
func (r *GongYiDraf) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=?", r.ShafaId).Delete(r).Error
	return err
}

func (r *GongYiDraf) Update(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Save(r).Error
	return err
}

func (r *GongYiDraf) DeleteByTypes(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=? and types = ?", r.ShafaId, r.Types).Delete(r).Error
	return err
}

func (r *GongYiDraf) DeleteWithTrans(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=? and types = ? and trans_id =? ", r.ShafaId, r.Types, r.TransId).Delete(r).Error
	return err
}

func (r *GongYiDraf) DeleteByFenweiName(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=? and fen_wei_name = ?", r.ShafaId, r.FenWeiName).Delete(r).Error
	return err
}

func (r *GongYiDraf) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("cp_type=?", cpType).Find(&r).Error
	return err
}

func (rList *GongYiDrafList) GetListPage(shafaId, types string, d *gorm.DB) (err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(gongyidrafTableName)
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

func (rList *GongYiDrafList) GetListPageWithTrans(shafaId, types, transId string, d *gorm.DB) (err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(gongyidrafTableName)
	if shafaId != "" {
		//说明查询条件不为空
		query = query.Where("shafa_id  =? ", shafaId)
	}

	if types != "" {
		//说明查询条件不为空
		query = query.Where("types  =? ", types)
	}
	if transId != "" {
		//说明查询条件不为空
		query = query.Where("trans_id  =? ", transId)
	}

	err = query.Order("gong_yi_name").Find(&rList).Error
	return

}

func (r *GongYiDrafList) GetBySoFaCode(d *gorm.DB, sofaCode string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=?", sofaCode).Find(&r).Error
	return err
}

func (r *GongYiDrafList) GetBySoFaCodeDraf(d *gorm.DB, sofaCode, transID string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(gongyidrafTableName).Where("shafa_id=? and  trans_id  =?", sofaCode, transID).Find(&r).Error
	return err
}
