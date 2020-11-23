package server

import (
	"encoding/base64"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_ca/src"
	"simple_ca/src/dao"
	"simple_ca/src/definition"
	"simple_ca/src/message"
	"simple_ca/src/tools"
	"strconv"
)

func CaRequestLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.CaCodeSignatureRequestReq)
	resp := message.CaRequestResp{}
	user, ok := ctx.Get("user")
	if !ok {
		resp.Header = ginTools.UnauthorizedRespHeader
		return resp
	}
	u := user.(*dao.User)

	// base64 解码
	msgSplit, err := base64.StdEncoding.DecodeString(request.CSRID)
	if err != nil {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	// DES 解密
	csrIDString, ok := tools.DecryptWithDES(msgSplit, src.GetSetting().Secret.ResponseSecret)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	// 转换成 int
	csrID, err := strconv.Atoi(csrIDString)
	if err != nil {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	csr, ok := dao.GetCRSByID(uint(csrID))

	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	if csr.UserID != u.ID {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusForbidden,
			Msg:  "你无权修改此项资源",
		}
		return resp
	}

	if csr.State != definition.CRSStateInit {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "请勿重复提交",
		}
		return resp
	}

	_, ok = dao.AddPublicKeyForRequest(csr, request.PublicKey, u.ID)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

func CaCsrLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.CaCsrReq)
	resp := message.CaCsrResp{}
	user, ok := ctx.Get("user")
	if !ok {
		resp.Header = ginTools.UnauthorizedRespHeader
		return resp
	}
	u := user.(*dao.User)
	// 存库
	newCSR := &dao.CARequest{
		UserID:               u.ID,
		State:                uint(1),
		PublicKey:            "",
		Country:              request.Country,
		Province:             request.Province,
		Locality:             request.Locality,
		Organization:         request.Organization,
		OrganizationUnitName: request.OrganizationalUnit,
		CommonName:           request.CommonName,
		EmailAddress:         request.EmailAddress,
	}
	newCSR, ok = dao.CreateNewCRS(newCSR)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	encryptID, ok := tools.EncryptWithDES(strconv.Itoa(int(newCSR.ID)), src.GetSetting().Secret.ResponseSecret)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	base64ID := base64.StdEncoding.EncodeToString(encryptID)
	resp.CSRID = base64ID
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
