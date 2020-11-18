package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/server"
)

func AuthRegister(rg *gin.RouterGroup) {
	rg.POST("/login", authLoginRoutes()...)
	rg.POST("/register", authRegisterRoutes()...)
}

func authLoginRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ginTools.EasyHandler(check.AuthLoginCheck,
			server.AuthLoginLogic, message.AuthLoginReq{}),
	}
}
func authRegisterRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ginTools.EasyHandler(check.AuthRegisterCheck,
			server.AuthRegisterLogic, message.AuthRegisterReq{}),
	}
}
