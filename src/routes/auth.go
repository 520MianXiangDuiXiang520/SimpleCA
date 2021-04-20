package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
	middlewareTools "github.com/520MianXiangDuiXiang520/ginUtils/middleware"
	"github.com/gin-gonic/gin"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/middleware"
	"simple_ca/src/server"
)

func AuthRegister(rg *gin.RouterGroup) {
	rg.POST("/login", authLoginRoutes()...)
	rg.POST("/register", authRegisterRoutes()...)
	rg.POST("/logout", authLogoutRoutes()...)
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

func authLogoutRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.AuthLogoutCheck,
			server.AuthLogoutLogic, message.AuthLogoutReq{}),
	}
}
