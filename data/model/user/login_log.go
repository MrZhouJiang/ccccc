package model

import (
	"ccccc/db"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	login_log_tableName = "login_log"
)

//登陆记录表
type LoginLog struct {
	Id      int       `json:"id"`
	LogTime time.Time `json:"log_time"`
	UserId  string    `json:"user_id"`
	//登陆类型  1 用户名密码 2 电话验证码 3 邮箱密码
	LoginType int `json:"login_type"`
	// 登陆渠道 1 默认web
	Channel int `json:"channel"`
}

type LoginLogList []LoginLog

func (r *LoginLog) Create(d *gorm.DB) error {
	r.LogTime = time.Now()
	if d == nil {
		d = db.BaseDB
	}
	err := d.Table(login_log_tableName).Create(r).Error
	return err
}
func (rList *LoginLogList) GetAllByUserId(userId string) error {
	if userId != "" {
		err := db.BaseDB.Table(login_log_tableName).Where("user_id= ?", userId).Order("log_time desc").Find(&rList).Error
		return err
	}
	err := db.BaseDB.Table(login_log_tableName).Order("log_time desc").Find(&rList).Error
	return err
}
