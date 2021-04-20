package check

import (
	"errors"
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/message"
)

func AuthLoginCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	request := req.(*message.AuthLoginReq)
	if request.Password == "" || request.Username == "" {
		return ginTools.ParamErrorRespHeader, errors.New("paramError")
	}
	return http.StatusOK, nil
}
func AuthRegisterCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	request := req.(*message.AuthRegisterReq)
	if request.Password == "" || request.Username == "" {
		return ginTools.ParamErrorRespHeader, errors.New("paramError")
	}
	return http.StatusOK, nil
}

func AuthLogoutCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}
