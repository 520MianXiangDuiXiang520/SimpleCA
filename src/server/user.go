package server

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"log"
	"simple_ca/src/message"
)

func UserCerLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.UserCerReq)
	resp := message.UserCerResp{}
	// TODO:...
	log.Println(request)
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
