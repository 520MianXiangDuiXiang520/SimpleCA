package server

import (
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
	"simple_ca/src/dao/utils"
	"simple_ca/src/message"
	"simple_ca/src/tools"
)

// 登录逻辑
func AuthLoginLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuthLoginReq)
	resp := message.AuthLoginResp{}
	pwd := tools.HashBySHA256([]string{request.Password})
	user, ok := dao.HasUser(request.Username, pwd)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	// md5 获取 token
	token := tools.HashByMD5([]string{
		pwd, fmt.Sprintf("%v", request),
	})
	// 写入 token
	if !dao.InsertToken(user, token) {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Token = token
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

// 注册逻辑
func AuthRegisterLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuthRegisterReq)
	resp := message.AuthRegisterResp{}
	pwd := tools.HashBySHA256([]string{request.Password})
	newUser := &utils.User{
		Username: request.Username,
		Password: pwd,
		Email:    request.Email,
	}
	if !dao.InsertUser(newUser) {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}