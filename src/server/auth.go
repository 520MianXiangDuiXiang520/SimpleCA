package server

import (
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src/dao"
	"simple_ca/src/message"
	"simple_ca/src/tools"
	"time"
)

// 登录逻辑
func AuthLoginLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuthLoginReq)
	resp := message.AuthLoginResp{}
	pwd := tools.HashBySHA256([]string{request.Password})
	user, ok := dao.HasUserByUP(request.Username, pwd)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	// md5 获取 token
	token := tools.HashByMD5([]string{
		pwd, fmt.Sprintf("%v%v", request, time.Now().Nanosecond()),
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
	if _, ok := dao.GetUserByName(request.Username); ok {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "用户名已存在！",
		}
		return resp
	}
	pwd := tools.HashBySHA256([]string{request.Password})
	newUser := &dao.User{
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

func AuthLogoutLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	resp := message.AuthLogoutResp{}
	user, ok := ctx.Get("user")
	if !ok {
		resp.Header = ginTools.UnauthorizedRespHeader
		return resp
	}
	u := user.(*dao.User)
	if !dao.DeleteTokenByUserID(u.ID) {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
