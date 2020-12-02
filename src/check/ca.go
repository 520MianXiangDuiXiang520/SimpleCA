package check

import (
	"errors"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/message"
	"simple_ca/src/tools"
)

func checkPublicKey(str string) bool {
	_, ok := tools.DecodeRSAPublicKey([]byte(str))
	return ok
}

func CaUploadPKCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	r := req.(*message.CaUploadPKReq)
	if !checkPublicKey(r.PublicKey) {
		resp := message.CaCsrResp{
			Header: ginTools.ParamErrorRespHeader,
		}
		return resp, errors.New("")
	}
	return http.StatusOK, nil
}

func CaCodeSignCsrCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}

func CaRevokeCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}

func CaFileCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	// request := req.(*message.CaCSRFileReq)

	return http.StatusOK, nil
}

func CaUpdateCrlCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}
func CaSslCsrCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}
