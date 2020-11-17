package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/server"
)

func ArticleRegister(rg *gin.RouterGroup) {
	rg.POST("/request", caRequestRoutes()...)
}

func caRequestRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ginTools.EasyHandler(check.CaRequestCheck,
			server.CaRequestLogic, message.CaRequestReq{}),
	}
}
