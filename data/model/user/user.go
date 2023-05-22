package model

import (
	"ccccc/db"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	user_tableName = "user_info"
)

//用户信息表
type UserInfo struct {
	Id int `json:"id"`
	//用户唯一ID
	UserId string `json:"user_id"`
	//昵称
	Name string `json:"name"`
	//头像
	Portrait string `json:"portrait"`
	//密文密码
	PassWord string `json:"pass_word"`
	//电话密文
	Phone string `json:"phone"`
	//油箱
	Email string `json:"email"`
	// 年龄
	Age int `json:"age"`
	//体重  单位 克
	Weight int `json:"weight"`
	//身份证号 密文
	IdCard string `json:"id_card"`
	//地址
	Address string `json:"address"`
	//创建时间
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	//备注
	DescInfo   string `json:"desc_info"`
	CreateUser string `json:"create_user"`
}

type UserInfoList []UserInfo

func (r *UserInfo) Create(d *gorm.DB) error {
	r.CreateTime = time.Now()
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(user_tableName).Create(r).Error
	return err
}
func (r *UserInfo) Delete(d *gorm.DB) error {
	r.CreateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(user_tableName).Where("user_id=?", r.UserId).Delete(r).Error
	return err
}

func (r *UserInfo) Get(userId, nick string) error {
	if userId != "" {
		err := db.BaseDB.Table(user_tableName).Where("user_id= ?", userId).Find(&r).Error
		return err
	}
	if nick != "" {
		err := db.BaseDB.Table(user_tableName).Where("name= ?", nick).Find(&r).Error
		return err
	}
	return nil
}

func (r *UserInfo) Update(d *gorm.DB) error {
	r.UpdateTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(user_tableName).Save(r).Error
	return err
}

func (list *UserInfoList) GetList() error {
	//剔除管理员账号

	return nil
}

func (rList *UserInfoList) GetListPage(userName string, offset, limit int, d *gorm.DB) (total int, err error) {
	if d == nil {
		d = db.BaseDB
	}
	query := d.Table(user_tableName)
	if userName != "" {
		//说明查询条件不为空
		query = query.Where("name like ? ", fmt.Sprintf("%%%s%%", userName))
	}

	query.Count(&total)
	err = query.Order("id desc").Offset(offset).Limit(limit).Find(&rList).Error
	return

}
