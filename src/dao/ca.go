package dao

import (
	"errors"
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/jinzhu/gorm"
	"simple_ca/src/definition"
)

func CreateNewCRS(request *CARequest) (res *CARequest, ok bool) {
	res, err := insertCRS(request)
	if err != nil {
		msg := fmt.Sprintf("Fail to insert csr: %v", request)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return res, true
}

func insertCRS(request *CARequest) (res *CARequest, err error) {
	if _, ok := HasUserByID(request.UserID); !ok {
		return nil, errors.New(fmt.Sprintf("no this user(id = %d)", request.UserID))
	}
	err = daoUtils.GetDB().Create(request).Error
	return request, err
}

func HasCRSByID(id uint) (ok bool) {
	crs, err := selectCRSByID(id)
	if err != nil {
		msg := fmt.Sprintf("Fail to select crs by id(%d)", id)
		utils.ExceptionLog(err, msg)
		return false
	}
	return crs.ID == id
}

func GetCRSByID(id uint) (*CARequest, bool) {
	csr, err := selectCRSByID(id)
	if err != nil || csr.ID != id {
		msg := fmt.Sprintf("Fail to select crs by id(%d)", id)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return csr, true
}

func selectCRSByID(id uint) (crs *CARequest, err error) {
	crs = &CARequest{}
	err = daoUtils.GetDB().Where("id = ?", id).First(crs).Error
	return crs, err
}

func AddPublicKeyForRequest(csr *CARequest, pk string, uid uint) (*CARequest, bool) {
	csr.PublicKey = pk
	csr.State = definition.CRSStateAuditing
	err := updateCRSByID(daoUtils.GetDB(), csr, csr.ID)
	if err != nil {
		msg := fmt.Sprintf("Fail to update csr, csr = %v", csr)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return csr, true
}

func updateCRSByID(db *gorm.DB, crs *CARequest, id uint) (err error) {
	err = db.Model(&CARequest{}).Where("id = ?",
		id).Update(crs).Error
	return err
}

func updateCRSStateByID(db *gorm.DB, state, id uint) (err error) {
	err = db.Model(&CARequest{}).Where("id = ?",
		id).Updates(map[string]interface{}{"state": state}).Error
	return err
}
