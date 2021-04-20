package dao

import (
	"errors"
	"fmt"
	dao_tools "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"github.com/jinzhu/gorm"
	"simple_ca/src/definition"
	"simple_ca/src/tools"
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
		tools.ExceptionLog(err,
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

func SetCSRState(crs *CARequest, state uint) (*CARequest, bool) {
	var err error
	oldState := crs.State
	db := dao_tools.GetDB()
	switch state {
	case definition.CRSStateInit:
		crs.State = definition.CRSStateInit
		err = updateCRSByID(db, crs, crs.ID)
	case definition.CRSStateAuditing:
		crs.State = definition.CRSStateAuditing
		err = updateCRSByID(db, crs, crs.ID)
	case definition.CRSStatePass:
		crs.State = definition.CRSStatePass
		err = updateCRSByID(db, crs, crs.ID)
	case definition.CRSStateUnPass:
		crs.State = definition.CRSStateUnPass
		err = updateCRSByID(db, crs, crs.ID)
	default:
		err = errors.New("UndefinedStatus")
		return nil, false
	}
	if err != nil {
		tools.ExceptionLog(err,
			fmt.Sprintf("Fail to update csr state from %d to %d", oldState, state))
		return nil, false
	}
	return crs, true
}

func CreateNewCertificate(c *Certificate) (*Certificate, bool) {
	r, err := insertNewCertificate(dao_tools.GetDB(), c)
	if err != nil || r.ID == 0 {
		tools.ExceptionLog(err,
			fmt.Sprintf("Fail to insert Certificate： %v", c))
		return nil, false
	}
	return c, true
}

// 插入一个新的证书
func insertNewCertificate(db *gorm.DB, c *Certificate) (*Certificate, error) {
	err := db.Create(c).Error
	return c, err
}
