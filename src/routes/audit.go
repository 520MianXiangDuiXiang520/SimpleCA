package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/server"
)

func AuditRegister(rg *gin.RouterGroup) {
	rg.POST("/list", auditListRoutes()...)
}

func auditListRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ginTools.EasyHandler(check.AuditListCheck,
			server.AuditListLogic, message.AuditListReq{}),
	}
}
