package server

import (
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
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
	fmt.Println("CSRID", request.CSRID, request.PublicKey)
	csrIDString := tools.DecryptWithDES(request.CSRID)
	csrID, _ := strconv.Atoi(csrIDString)
	_, ok = dao.AddPublicKeyForRequest(uint(csrID), request.PublicKey, u.ID)
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
	encryptID := tools.EncryptWithDES(strconv.Itoa(int(newCSR.ID)))
	resp.CSRID = encryptID
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
