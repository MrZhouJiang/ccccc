package router

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	// 单个课程
	HomeRouter(engine)

}
