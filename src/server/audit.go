package server

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
	"simple_ca/src/definition"
	"simple_ca/src/message"
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
