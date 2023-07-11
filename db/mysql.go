package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var mysql_db *gorm.DB

//const DSN = "root:codoon20140312@tcp(120.26.10.118:3306)/tt_test?charset=utf8mb4&parseTime=true&loc=Local"

const DSN = "root:admin123@tcp(127.0.0.1:3306)/cp?charset=utf8mb4&parseTime=true&loc=Local"

//指定驱动
const DRIVER = "mysql"

func InitMySqlDb() *gorm.DB {
	fmt.Println(" hello init mysql")
	var err error
	mysql_db, err = gorm.Open(DRIVER, DSN)

	if err != nil {
		fmt.Println(" hello init InitMySqlDb err " + err.Error())
		log.Printf("init mysql_db error :%v\n", err)

		panic(err)
	} else {
		fmt.Println("init mysql success")
	}
	return mysql_db

}
