package server

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	utils "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

func CaCrlLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	request := req.(*message.CaCrlReq)
	resp := message.CaCrlResp{}
	// 获取用户身份
	user, ok := ctx.Get("user")
	if !ok {
		resp.Header = ginTools.UnauthorizedRespHeader
		return resp
	}
	u := user.(*dao.User)

	// 检查证书
	cer, ok := dao.GetCertificateByID(request.SerialNumber)
	if !ok {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "证书不存在",
		}
		return resp
	}

	if cer.UserID != u.ID {
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusForbidden,
			Msg:  "拒绝服务！",
		}
		return resp
	}

	// crl 信息落库
	cTime := time.Now().Unix()
	_, err := dao.CreateNewCRL(cer.RequestID, request.SerialNumber, cTime)
	if err != nil {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	// 更新 crl 文件
	ok = updateCRLFile()
	if !ok {
		resp.Header = ginTools.SystemErrorRespHeader
		return resp
	}
	resp.Header = ginTools.SuccessRespHeader
	return resp
}

// 更新 CRL 文件
func updateCRLFile() bool {
	if time.Now().Unix() < src.GetNextUpdateCRLTime() {
		return true
	}
	crlList, err := dao.GetAllCRL()
	if err != nil {
		return false
	}
	rcList := make([]pkix.RevokedCertificate, len(crlList))
	for i, v := range crlList {
		rcList[i] = pkix.RevokedCertificate{
			SerialNumber:   big.NewInt(int64(v.CertificateID)),
			RevocationTime: time.Unix(v.InputTime, 0),
		}
	}
	rootCer, rootPK := src.GetCARootCer()
	n := time.Now()
	l := n.Add(time.Hour * 24)
	ok := tools.CreateNewCRL(&rootCer, &rootPK, rcList, n, l, src.GetSetting().CRLSetting.CRLFileName)
	src.SetNextUpdateCRLTime(n.Unix())
	return ok
}

func CaFileLogic(ctx *gin.Context, req ginTools.BaseReqInter) ginTools.BaseRespInter {
	// request := req.(*message.CaFileReq)
	resp := message.CaFileResp{}
	fileHeader, err := ctx.FormFile(src.GetSetting().CSRFileKey)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to read %s", src.GetSetting().CSRFileKey))
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "请选择要上传的文件",
		}
		return resp
	}
	if fileHeader.Size <= 0 {
		utils.ExceptionLog(err, fmt.Sprintf("The size of %s is %d", fileHeader.Filename, fileHeader.Size))
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "空文件",
		}
		return resp
	}
	file, err := fileHeader.Open()
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to open %s", fileHeader.Filename))
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusInternalServerError,
			Msg:  "无法打开文件",
		}
		return resp
	}
	fileBytes, _ := ioutil.ReadAll(file)
	block, rest := pem.Decode(fileBytes)

	if block == nil || len(rest) > 0 {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to parse %s", fileHeader.Filename))
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "无法解析文件",
		}
		return resp
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to parse %s", fileHeader.Filename))
		resp.Header = ginTools.BaseRespHeader{
			Code: http.StatusBadRequest,
			Msg:  "无法解析文件",
		}
		return resp
	}
	resp.Country = getFirstBySplit(csr.Subject.Country)
	resp.Province = getFirstBySplit(csr.Subject.Province)
	resp.Locality = getFirstBySplit(csr.Subject.Locality)
	resp.Organization = getFirstBySplit(csr.Subject.Organization)
	resp.OrganizationalUnit = getFirstBySplit(csr.Subject.OrganizationalUnit)
	resp.CommonName = csr.Subject.CommonName
	resp.EmailAddress = getFirstBySplit(csr.EmailAddresses)

	resp.Header = ginTools.SuccessRespHeader
	return resp
}

func getFirstBySplit(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}
