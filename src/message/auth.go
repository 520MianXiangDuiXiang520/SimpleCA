package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
)

type AuthLoginResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
	Token  string                  `json:"token"`
}

type AuthLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *AuthLoginReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type AuthRegisterResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type AuthRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r *AuthRegisterReq) JSON(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(&r)
	return err
}
