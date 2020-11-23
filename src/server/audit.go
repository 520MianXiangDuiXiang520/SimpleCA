package server

import (
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	utils "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"simple_ca/src"
	"simple_ca/src/dao"
	"simple_ca/src/definition"
	"simple_ca/src/message"
	"simple_ca/src/tools"
	"strconv"
	"time"
)

func AuditListLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	resp := message.AuditListResp{}
	list, ok := dao.GetCRSsByState(definition.CRSStateAuditing)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.CRSList = make([]message.CRSPublicKey, len(list))
	for i, v := range list {
		k := message.CRSPublicKey{}
		k.ID = v.ID
		k.PublicKey = v.PublicKey
		k.CommonName = v.CommonName
		k.Organization = v.Organization
		k.Locality = v.Locality
		k.Province = v.Province
		k.EmailAddress = v.EmailAddress
		k.OrganizationalUnit = v.OrganizationUnitName
		k.Country = v.Country
		resp.CRSList[i] = k
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

func AuditPassLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuditPassReq)
	resp := message.AuditPassResp{}
	// 管理员身份验证
	csrIDBytes, err := base64.StdEncoding.DecodeString(request.CSRID)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Base64 decoding failed， input = %s", request.CSRID))
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	csrIDStr, ok := tools.DecryptWithDES(csrIDBytes, src.GetSetting().Secret.ResponseSecret)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	csrID, err := strconv.Atoi(csrIDStr)
	if err != nil {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	csr, ok := dao.GetCRSByID(uint(csrID))
	if csr.State != definition.CRSStateAuditing {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	// 修改 CSR 状态
	csr, ok = dao.SetCSRState(csr, definition.CRSStateUnPass)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	// 落库
	// long := src.GetSetting().Secret.CertificateEffectiveTime
	notBefore := time.Now()
	notAfter := time.Now().Add(time.Hour * 24 * 365)
	expireTime := time.Now().Unix() + definition.WrongOneYear
	c, ok := dao.CreateNewCertificate(&dao.Certificate{
		State:      definition.CertificateStateUsing,
		ExpireTime: expireTime,
		UserID:     csr.UserID,
		RequestID:  csr.ID,
	})
	// 生成证书
	rootCer, rootPK := src.GetCARootCer()
	ok = tools.CreateNewCertificate(&rootCer, big.NewInt(int64(int(c.ID))), pkix.Name{
		Country:            []string{csr.Country},
		Province:           []string{csr.Province},
		Locality:           []string{csr.Locality},
		Organization:       []string{csr.Organization},
		OrganizationalUnit: []string{csr.OrganizationUnitName},
		CommonName:         csr.CommonName,
	}, csr.PublicKey, &rootPK, notBefore, notAfter,
		fmt.Sprintf("../cers/junebao_%v.cer", time.Now().Nanosecond()))
	if !ok {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusInternalServerError,
			Msg:  "证书生成失败！",
		}
		return resp
	}
	// 邮件通知用户
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
