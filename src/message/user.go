package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"simple_ca/src/definition"
)

type UserCerResp struct {
	Header       ginTools.BaseRespHeader                  `json:"header"`
	Certificates []definition.CertificateFullAmountFields `json:"certificates"`
}

type UserCerReq struct {
}

func (r *UserCerReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
