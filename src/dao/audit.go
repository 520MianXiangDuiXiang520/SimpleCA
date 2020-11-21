package dao

import (
	"errors"
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	logTools "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"simple_ca/src/definition"
)

func GetCRSsByState(state uint) (res []CARequest, ok bool) {
	var err error
	switch state {
	case definition.CRSStateInit:
		res, err = selectCRSsByState(definition.CRSStateInit)
	case definition.CRSStateAuditing:
		res, err = selectCRSsByState(definition.CRSStateAuditing)
	case definition.CRSStatePass:
		res, err = selectCRSsByState(definition.CRSStatePass)
	case definition.CRSStateUnPass:
		res, err = selectCRSsByState(definition.CRSStateUnPass)
	default:
		err = errors.New("UndefinedStatus")

	}
	if err != nil {
		logTools.ExceptionLog(err,
			fmt.Sprintf("Failed to query CRS list by status （%d）", state))
		return nil, false
	}
	return res, true
}

func selectCRSsByState(state uint) ([]CARequest, error) {
	res := make([]CARequest, 0)
	err := dao_tools.GetDB().Model(&CARequest{}).Where("state = ?", state).Find(&res).Error
	return res, err
}
