package message

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
)

type CaRequestResp struct {
	Header ginTools.BaseRespHeader `json:"header"`
}

type CaRequestReq struct {
	PublicKey    string `json:"public_key"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Organization string `json:"organization"`
	UnitName     string `json:"unit_name"`
	CommonName   string `json:"common_name"`
	EmailAddress string `json:"email_address"`
}

func (r CaRequestReq) JSON(ctx *gin.Context) error {
	return ctx.ShouldBindJSON(&r)
}
