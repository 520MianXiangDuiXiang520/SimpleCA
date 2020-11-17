package server

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"log"
	"simple_ca/src/message"
)

func CaRequestLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.CaRequestReq)
	resp := message.CaRequestResp{}
	// TODO:...
	log.Println(request)
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
