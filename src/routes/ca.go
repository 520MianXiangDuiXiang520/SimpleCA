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

func CARegister(rg *gin.RouterGroup) {
	rg.POST("/request", caRequestRoutes()...)
	rg.POST("/csr", caCsrRoutes()...)
	rg.POST("/crl", caCrlRoutes()...)

}

func caRequestRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaRequestCheck,
			server.CaRequestLogic, message.CaCodeSignatureRequestReq{}),
	}
}

func caCsrRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaCsrCheck,
			server.CaCsrLogic, message.CaCsrReq{}),
	}
}

func caCrlRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaCrlCheck,
			server.CaCrlLogic, message.CaCrlReq{}),
	}
}
