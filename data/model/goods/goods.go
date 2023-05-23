package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

//物料信息 基本表
//同步时间 每天清晨

type Goods struct {
	Id         int    `json:"id"`
	CpCode     string `json:"cp_code"`
	CpName     string `json:"cp_name"`
	CpDesc     string `json:"cp_desc"`
	CpType     string `json:"cp_type"`
	CpTypeCode string `json:"cp_type_code"`
	//规格(成品有)
	CpGuiGe string `json:"cp_gui_ge"`
	//主计量单位
	CpMainUnit string `json:"cp_main_unit"`
	//辅计量单位
	FuZhuUnit string `json:"fu_zhu_unit"`
	//产品尺寸 3*5 （cm） 注意 都是里面单位 3 表示3 cm   （3*5*2）
	MainSize string `json:"main_size"`
	//主系数
	MainXiShu float64 `json:"main_xi_shu"`
	//辅助系数
	FuZhuXiShu float64   `json:"fu_zhu_xi_shu"`
	Price      float64   `json:"price"`
	ShunHao    string    `json:"shun_hao"`
	ChangeP    float64   `json:"change_p"`
	LoadTime   time.Time `json:"load_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	//转换率
	DomH string `json:"dom_h"`
	//转换配置ID
	DomId int `json:"dom_id"`
	//通用名称
	TyName string `json:"ty_name"`
}

const tabel_name = "goods"

type GoodsList []Goods

func (r *Goods) Create(d *gorm.DB) error {
	r.LoadTime = time.Now()
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Create(r).Error
	return err
}
func (r *Goods) Delete(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Where("cp_code=?", r.CpCode).Delete(r).Error
	return err
}

func (r *Goods) GetById(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Where("id=?", r.Id).Find(&r).Error
	return err
}
func (r *Goods) Save(updateData map[string]interface{}, d *gorm.DB) error {

	if d == nil {
		d = db.BaseDB
	}
	updateData["update_time"] = time.Now()
	err := d.Table(tabel_name).Where("id = ?", r.Id).Updates(updateData).Error
	return err
}
func (r *Goods) GetByCpCode(d *gorm.DB) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Where("cp_code=?", r.CpCode).Find(&r).Error
	return err
}

func (r *GoodsList) GetByType(d *gorm.DB, cpType string) error {
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Where("cp_type=?", cpType).Find(&r).Error
	return err
}

func (rList *GoodsList) GetListPage(name, code, types string, offset, limit int, d *gorm.DB) (total int, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(tabel_name)
	if name != "" {
		//说明查询条件不为空
		query = query.Where("cp_name like ? ", fmt.Sprintf("%%%s%%", name))
	}

	if code != "" {
		//说明查询条件不为空
		query = query.Where("cp_code like ? ", fmt.Sprintf("%%%s%%", code))
	}
	if types != "" {
		//说明查询条件不为空
		query = query.Where("cp_type_code = ? ", types)
	}

	query.Count(&total)
	err = query.Order("id desc").Offset(offset).Limit(limit).Find(&rList).Error
	return

}

func (r *Goods) Get(name, code string) error {
	if name != "" {
		err := db.BaseDB.Table(tabel_name).Where("cp_name= ?", name).Find(&r).Error
		return err
	}
	if code != "" {
		err := db.BaseDB.Table(tabel_name).Where("cp_code= ?", code).Find(&r).Error
		return err
	}
	return nil
}

func (r *Goods) Update(d *gorm.DB) error {
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(tabel_name).Save(r).Error
	return err
}

type AllGoodsDesc struct {
	CpCode       string  `json:"cp_code"`
	CpName       string  `json:"cp_name"`
	CpGuiGe      string  `json:"cp_gui_ge"`
	CpMainUnit   string  `json:"cp_main_unit"`
	FuZhuUnit    string  `json:"fu_zhu_unit"`
	CpMainUnitId int     `json:"cp_main_unit_id"`
	FuZhuUnitId  int     `json:"fu_zhu_unit_id"`
	MainXiShu    float64 `json:"main_xi_shu"`
	FuZhuXiShu   float64 `json:"fu_zhu_xi_shu"`
	//产品尺寸 3*5 （cm） 注意 都是里面单位 3 表示3 cm   （3*5*2）
	MainSize   string  `json:"main_size"`
	Price      float64 `json:"price"`
	ShunHao    string  `json:"shun_hao"`
	ChangeP    float64 `json:"change_p"`
	CpTypeCode string  `json:"cp_type_code"`
	CpType     string  `json:"cp_type"`
	CpDesc     string  `json:"cp_desc"`
}

func (rList *GoodsList) GetAllGoodsList(d *gorm.DB) (list []AllGoodsDesc, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(tabel_name)
	query = query.Select("cp_code , cp_name,cp_gui_ge,cp_main_unit,fu_zhu_unit,main_xi_shu,fu_zhu_xi_shu,main_size,price,shun_hao,change_p,cp_type_code,cp_desc")
	err = query.Order("id desc").Find(&list).Error
	return

}

func (rList *Goods) GetGoodsById(cpCode string, d *gorm.DB) (info AllGoodsDesc, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(tabel_name)
	query = query.Select("cp_code , cp_name,cp_gui_ge,cp_main_unit,fu_zhu_unit,main_xi_shu,fu_zhu_xi_shu,main_size,price,shun_hao,change_p,cp_type_code,cp_desc")
	err = query.Where("cp_code = ?", cpCode).First(&info).Error
	return

}
