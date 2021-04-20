package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
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
	Username string `json:"username" check:"len: [5, 16]; not null"`
	Password string `json:"password" check:"len: [5, 18]; not null"`
	Email    string `json:"email"    check:"email; not null"`
}

func (r *AuthRegisterReq) JSON(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(&r)
	return err
}

type AuthLogoutResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type AuthLogoutReq struct {
}

func (r *AuthLogoutReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
