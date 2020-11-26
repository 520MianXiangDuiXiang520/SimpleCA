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
	// 管理员身份验证
	csr, ok := checkCSRID(request.CSRID)
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
	cerFileName := fmt.Sprintf("%s/%d_%v_%v.cer",
		src.GetSetting().Secret.UserCerPath, user.ID,
		user.Username, time.Now().Nanosecond())
	// 获取 CA 根证书和私钥
	rootCer, rootPK := src.GetCARootCer()
	// 根据 CA 根证书和私钥为用户签发证书
	ok = tools.CreateNewCertificate(&rootCer, big.NewInt(int64(int(c.ID))), pkix.Name{
		Country:            []string{csr.Country},
		Province:           []string{csr.Province},
		Locality:           []string{csr.Locality},
		Organization:       []string{csr.Organization},
		OrganizationalUnit: []string{csr.OrganizationUnitName},
		CommonName:         csr.CommonName,
	}, csr.PublicKey, &rootPK, notBefore, notAfter,
		[]string{src.GetSetting().CRLSetting.CRLDistributionPoint}, cerFileName)
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
	csr, ok := checkCSRID(request.CSRID)
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
		Subject: "证书申请通过通知",
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
