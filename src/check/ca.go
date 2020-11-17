package check

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/message"
)

func CaRequestCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	_ = req.(*message.CaRequestReq)

	return http.StatusOK, nil
}
