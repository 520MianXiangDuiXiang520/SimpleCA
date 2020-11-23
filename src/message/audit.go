package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/definition"
)

type CRSPublicKey struct {
	ID        uint   `json:"id"`
	PublicKey string `json:"public_key"`
	definition.CertificateSigningRequest
}

type AuditListResp struct {
	Header  ginTools.BaseRespHeader `json:"header"`
	CRSList []CRSPublicKey          `json:"crs_list"`
}

type AuditListReq struct {
}

func (r *AuditListReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}

type AuditPassResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type AuditPassReq struct {
	CSRID string `json:"csr_id"` // CSR ID 3DES 加密后 Base64 编码
}

func (r *AuditPassReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
