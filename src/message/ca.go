package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
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
type CaUploadPKReq struct {
	CSRID     string `json:"csrid"      check:"not null"`
	PublicKey string `json:"public_key" check:"not null"`
}

func (r *CaUploadPKReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaCsrResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
	CSRID  string                  `json:"csr_id"`
}

type CaCodeSignCsrReq struct {
	definition.CertificateSigningRequest
}

func (r *CaCodeSignCsrReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaRevokeResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type CaRevokeReq struct {
	SerialNumber uint `json:"serial_number" check:"not null"` // 序列号
}

func (r *CaRevokeReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaCSRFileResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
	definition.CertificateSigningRequest
	PublicKey string `json:"public_key"`
}

type CaCSRFileReq struct {
}

func (r *CaCSRFileReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaUpdateCrlResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type CaUpdateCrlReq struct {
}

func (r *CaUpdateCrlReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type CaSslCsrResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
	CSRID  string                  `json:"csr_id"`
}

type CaSslCsrReq struct {
	definition.CertificateSigningRequest
	DNSNames string `json:"dns_names" check:"not null"`
}

func (r *CaSslCsrReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
