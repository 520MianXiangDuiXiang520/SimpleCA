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

func CaRequestCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	r := req.(*message.CaCodeSignatureRequestReq)
	if !checkPublicKey(r.PublicKey) {
		resp := message.CaCsrResp{
			Header: ginTools.ParamErrorRespHeader,
		}
		return resp, errors.New("")
	}
	return http.StatusOK, nil
}

func CaCsrCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}

func CaCrlCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}

func CaFileCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	// request := req.(*message.CaFileReq)

	return http.StatusOK, nil
}
