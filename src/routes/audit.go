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

func AuditRegister(rg *gin.RouterGroup) {
	rg.POST("/list", auditListRoutes()...)
	rg.POST("/pass", auditPassRoutes()...)
	rg.POST("/unpass", auditUnPassRoutes()...)
}

func auditListRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		middlewareTools.Permiter(middleware.AdminPermit),
		ginTools.EasyHandler(check.AuditListCheck,
			server.AuditListLogic, message.AuditListReq{}),
	}
}

func auditPassRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		middlewareTools.Permiter(middleware.AdminPermit),
		ginTools.EasyHandler(check.AuditPassCheck,
			server.AuditPassLogic, message.AuditPassReq{}),
	}
}

func auditUnPassRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		middlewareTools.Permiter(middleware.AdminPermit),
		ginTools.EasyHandler(check.AuditUnPassCheck,
			server.AuditUnPassLogic, message.AuditUnPassReq{}),
	}
}
