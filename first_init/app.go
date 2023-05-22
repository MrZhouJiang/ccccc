package first_init

import "ccccc/db"

func InitDb() error {
	//app的一系列初始化操作
	db.InitDB()
	//defer db.CloseDb()
	return nil

}
