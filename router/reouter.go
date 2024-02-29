package router

import (
	"MURI-GO/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 实例化Router对象，可以使用该对象点出首字母大写的方法（挎包使用）
var Router router

// 定义router结构体
type router struct {
}

// 用于初始化路由
func (r *router) InitApiRouter(router *gin.Engine) {
	router.Use(logger.GinLogger(), logger.GinRecovery(true))
	router.
		//
		GET("/testapi", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "success",
			})
		})
}
