package server

import (
	"crypto/x509/pkix"
	"encoding/base64"
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/email_tools"
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
	"strings"
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
		switch v.Type {
		case definition.CertificateTypeCodeSign:
			k.TypeStr = "CodeSign"
		case definition.CertificateTypeSSL:
			k.TypeStr = "SSL"
		}
		resp.CRSList[i] = k
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

func checkCSRID(CSRID string) (*dao.CARequest, bool) {
	csrIDBytes, err := base64.StdEncoding.DecodeString(CSRID)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Base64 decoding failed， input = %s", CSRID))
		return nil, false
	}
	csrIDStr, ok := tools.DecryptWithDES(csrIDBytes, src.GetSetting().Secret.ResponseSecret)
	if !ok {
		return nil, false
	}
	csrID, err := strconv.Atoi(csrIDStr)
	if err != nil {
		return nil, false
	}
	csr, ok := dao.GetCRSByID(uint(csrID))
	if !ok {
		return nil, false
	}
	return csr, true
}

func AuditPassLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuditPassReq)
	resp := message.AuditPassResp{}
	// 删除解密步骤(2020/11/29)
	// csr, ok := checkCSRID(request.CSRID)
	csr, ok := dao.GetCRSByID(request.CSRID)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	if csr.State != definition.CRSStateAuditing {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	// 修改 CSR 状态
	csr, ok = dao.SetCSRState(csr, definition.CRSStatePass)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}

	user, ok := dao.HasUserByID(csr.UserID)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}

	notBefore := time.Now()
	// 证书有效时间为 365 天
	notAfter := time.Now().Add(time.Hour * 24 * 365)
	expireTime := time.Now().Unix() + definition.WrongOneYear

	// 落库存储
	c, ok := dao.CreateNewCertificate(&dao.Certificate{
		State:      definition.CertificateStateUsing,
		ExpireTime: expireTime,
		UserID:     csr.UserID,
		RequestID:  csr.ID,
	})
	// 生成证书
	cName := tools.GetCertificateFileName(c.ID, user.ID, user.Username)
	cerFileName := fmt.Sprintf("%s/%s",
		src.GetSetting().Secret.UserCerPath, cName)
	// 获取 CA 根证书和私钥
	rootCer, rootPK := src.GetCARootCer()

	subject := pkix.Name{
		Country:            []string{csr.Country},
		Province:           []string{csr.Province},
		Locality:           []string{csr.Locality},
		Organization:       []string{csr.Organization},
		OrganizationalUnit: []string{csr.OrganizationUnitName},
		CommonName:         csr.CommonName,
	}
	crlDP := []string{src.GetSetting().CRLSetting.CRLDistributionPoint}

	switch csr.Type {
	// 签发代码签名证书
	case definition.CertificateTypeCodeSign:
		ok = tools.CreateCodeSignCert(&rootCer, big.NewInt(int64(int(c.ID))), subject,
			csr.PublicKey, &rootPK, notBefore, notAfter, crlDP, cerFileName)
	// 签发 SSL 证书
	case definition.CertificateTypeSSL:
		dnsNames := strings.Split(csr.DnsNames, " ")
		ok = tools.CreateSSLCert(&rootCer, big.NewInt(int64(int(c.ID))), subject,
			csr.PublicKey, &rootPK, notBefore, notAfter, crlDP, dnsNames, cerFileName)
	default:
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}
	if !ok {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusInternalServerError,
			Msg:  "证书生成失败！",
		}
		return resp
	}
	// 邮件通知用户
	emailTemp := definition.CerSuccessTemp(map[string]string{
		"siteLink":    src.GetSetting().SiteLink,
		"username":    user.Username,
		"requestTime": csr.CreatedAt.Format("2006-01-02 15:04:05"),
		"time":        time.Now().Format("2006-01-02 15:04:05"),
	})
	err := email_tools.Send(&email_tools.EmailCTX{
		ToList: []email_tools.EmailUser{
			{Address: user.Email, Name: user.Username},
		},
		Subject: "证书申请通过通知",
		Body:    emailTemp,
		Path:    cerFileName,
	})
	if err != nil {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusInternalServerError,
			Msg:  "证书申请已通过，但颁发失败，请联系用户：" + user.Email,
		}
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

func AuditUnPassLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.AuditUnPassReq)
	resp := message.AuditUnPassResp{}
	// csr, ok := checkCSRID(request.CSRID)
	csr, ok := dao.GetCRSByID(request.CSRID)
	if !ok {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	if csr.State != definition.CRSStateAuditing {
		resp.Header = ginTools.ParamErrorRespHeader
		return resp
	}

	user, ok := dao.HasUserByID(csr.UserID)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}

	// 修改 CSR 状态
	csr, ok = dao.SetCSRState(csr, definition.CRSStateUnPass)
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}

	// 邮件通知
	emailTemp := definition.CerUnPassTemp(map[string]string{
		"siteLink":    src.GetSetting().SiteLink,
		"username":    user.Username,
		"requestTime": csr.CreatedAt.Format("2006-01-02 15:04:05"),
		"time":        time.Now().Format("2006-01-02 15:04:05"),
	})
	err := email_tools.Send(&email_tools.EmailCTX{
		ToList: []email_tools.EmailUser{
			{Address: user.Email, Name: user.Username},
		},
		Subject: "证书申请驳回通知",
		Body:    emailTemp,
	})
	if err != nil {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusInternalServerError,
			Msg:  "证书申请未通过，但邮件通知失败，请联系用户：" + user.Email,
		}
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}
