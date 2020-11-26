package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	middlewareTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/middleware"
	"github.com/gin-gonic/gin"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/middleware"
	"simple_ca/src/server"
)

func UserRegister(rg *gin.RouterGroup) {
	rg.POST("/cer", userCerRoutes()...)

}

func userCerRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.UserCerCheck,
			server.UserCerLogic, message.UserCerReq{}),
	}
}
