package check

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserCerCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}
