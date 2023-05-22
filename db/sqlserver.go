package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"log"
	"time"
)

var Sqlserver_db *gorm.DB

const DSN_sql = "sqlserver://zh2023:Admin123!@127.0.0.1/SFC"

//指定驱动
const DRIVER_sql_server = "mssql"

func InitSqlServerDb() *gorm.DB {
	fmt.Println(" hello init sqlserver113331444")

	var err error
	Sqlserver_db, err = gorm.Open(DRIVER_sql_server, DSN_sql)
	if err != nil {
		fmt.Println(" hello init sqlserver err " + err.Error())
		log.Printf("init sqlserver error :%v\n", err)
		panic(err)
	} else {
		fmt.Println("init sql server success")
	}
	Sqlserver_db.SingularTable(true)
	Sqlserver_db.LogMode(true)
	Sqlserver_db.DB().SetMaxIdleConns(10)
	Sqlserver_db.DB().SetMaxOpenConns(100)
	Sqlserver_db.DB().SetConnMaxLifetime(59 * time.Second)

	return Sqlserver_db

}
