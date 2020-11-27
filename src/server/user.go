package server

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
	"simple_ca/src/message"
)

func UserCerLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	resp := message.UserCerResp{}
	user, ok := ctx.Get("user")
	if !ok {
		resp.Header = ginTools.UnauthorizedRespHeader
		return resp
	}
	u := user.(*dao.User)
	res, ok := dao.GetCertificateFullAmountFieldsUser(u)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Certificates = res
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
