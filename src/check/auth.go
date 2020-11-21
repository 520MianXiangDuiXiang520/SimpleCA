package check

import (
	"errors"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"log"
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
		log.Println(request)
		return ginTools.ParamErrorRespHeader, errors.New("paramError")
	}
	return http.StatusOK, nil
}
