package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/definition"
)

type BaseCARequestReq struct {
	CSRID     string `json:"csrid"    check:"not null"`
	PublicKey string `json:"public_key" check:"not null"`
}

type CaRequestResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

// 代码签名请求头
type CaCodeSignatureRequestReq struct {
	CSRID     string `json:"csrid"      check:"not null"`
	PublicKey string `json:"public_key" check:"not null"`
}

func (r *CaCodeSignatureRequestReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaCsrResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
	CSRID  string                  `json:"csr_id"`
}

type CaCsrReq struct {
	definition.CertificateSigningRequest
}

func (r *CaCsrReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaCrlResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type CaCrlReq struct {
	SerialNumber uint `json:"serial_number" check:"not null"` // 序列号
}

func (r *CaCrlReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
