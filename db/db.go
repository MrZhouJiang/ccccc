package db

import (
	"github.com/jinzhu/gorm"
)

var BaseDB *gorm.DB
var SqlDB *gorm.DB

func InitDB() {
	db := InitMySqlDb()
	BaseDB = db

	/*	SDB := InitSqlServerDb()

		SqlDB = SDB*/

	/*	dddd := Test{}
		e1 := dddd.GetByType(nil, "")
		if e1 != nil {
			fmt.Printf("errrrrrrrrrrrr: %v\n", e1)
		} else {
			fmt.Printf("test sccess user: %v", dddd)
		}*/
}

func CloseDb() {
	defer BaseDB.Close() //退出前执行关闭
}
