package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
)

type UserCerResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type UserCerReq struct {
}

func (r *UserCerReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
