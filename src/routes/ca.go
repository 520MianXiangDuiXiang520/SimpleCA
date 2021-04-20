package routes

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
	middlewareTools "github.com/520MianXiangDuiXiang520/ginUtils/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/check"
	"simple_ca/src/message"
	"simple_ca/src/middleware"
	"simple_ca/src/server"
)

func CARegister(rg *gin.RouterGroup) {
	// 提交代码签名 CSR 信息
	rg.POST("/code_sign_csr", caCodeSignCsrRoutes()...)
	rg.POST("/ssl_csr", caSslCsrRoutes()...)

	// 单独提交公钥
	rg.POST("/upload_pk", caUploadPKRoutes()...)

	// 注销证书
	rg.POST("/revoke", caRevokeRoutes()...)

	// 上传 CSR 文件
	rg.POST("/csr_file", caCSRFileRoutes()...)

	// 手动更新 CRL 文件
	rg.POST("/update_crl", caUpdateCrlRoutes()...)

}

func caUploadPKRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaUploadPKCheck,
			server.CaUploadPKLogic, message.CaUploadPKReq{}),
	}
}

func caCodeSignCsrRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaCodeSignCsrCheck,
			server.CaCodeSignCsrLogic, message.CaCodeSignCsrReq{}),
	}
}

func caRevokeRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaRevokeCheck,
			server.CaRevokeLogic, message.CaRevokeReq{}),
	}
}

func caCSRFileRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		func(context *gin.Context) {
			resp := server.CaCSRFileLogic(context, &message.CaCSRFileReq{})
			context.Set("resp", resp)
			context.JSON(http.StatusOK, resp)
		},
	}
}

func caUpdateCrlRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		middlewareTools.Permiter(middleware.AdminPermit),
		ginTools.EasyHandler(check.CaUpdateCrlCheck,
			server.CaUpdateCrlLogic, message.CaUpdateCrlReq{}),
	}
}

func caSslCsrRoutes() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middlewareTools.Auth(middleware.TokenAuth),
		ginTools.EasyHandler(check.CaSslCsrCheck,
			server.CaSslCsrLogic, message.CaSslCsrReq{}),
	}
}
