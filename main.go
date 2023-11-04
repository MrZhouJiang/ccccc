package main

import (
	"ccccc/common"
	"ccccc/first_init"
	"ccccc/router"
	"ccccc/sysnc"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println(" hello wms")
	//初始化db
	err := first_init.InitDb()
	fmt.Println(" hello 1")
	if err != nil {
		log.Fatal(err)
	}
	//初始化 router
	engine := gin.New()
	engine.Use(common.GinLogger())
	engine.Use(common.TokenChecker())
	engine.Use(common.Core())
	router.Router(engine)

	common.InitBaseData()
	//sysnc.Task()
	sysnc.StartSyncCp()

	//sysnc.Task()
	//默认地址
	go sysnc.Task_base()
	engine.Run("192.168.202.5:8889")
	//engine.Run("127.0.0.1:8889")

}
