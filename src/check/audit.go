package check

import (
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/message"
)

func AuditListCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	request := req.(*message.AuditListReq)
	fmt.Println(request)
	return http.StatusOK, nil
}

func AuditPassCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	request := req.(*message.AuditPassReq)
	fmt.Println(request)
	return http.StatusOK, nil
}
