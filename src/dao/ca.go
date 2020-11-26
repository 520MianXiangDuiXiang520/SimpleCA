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
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to update csr state(%d) by id(%d)", state, id))
		return err
	}
	return err
}

func selectCertificateByID(db *gorm.DB, id uint) (*Certificate, error) {
	c := Certificate{}
	err := db.Where("id = ?", id).First(&c).Error
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to select Certificate by id(%d)", id))
		return nil, err
	}
	return &c, err
}

func GetCertificateByID(id uint) (Certificate, bool) {
	db := daoUtils.GetDB()
	c, err := selectCertificateByID(db, id)
	if err != nil {
		return Certificate{}, false
	}
	return *c, true
}

func insertNewCRL(db *gorm.DB, serial uint, expire int64) (*CRL, error) {
	crl := CRL{
		CertificateID: serial,
		InputTime:     expire,
	}
	err := db.Create(&crl).Error
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to insert crl: %v", crl))
		return nil, err
	}
	return &crl, nil
}

func updateCertificateStateByID(db *gorm.DB, state, id uint) error {
	err := db.Model(&Certificate{}).Where("id = ?",
		id).Updates(map[string]interface{}{"state": state}).Error
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to update Certificate state(%d) by id(%d)", state, id))
		return err
	}
	return err
}

// 生成一条新的 CRL 信息
func CreateNewCRL(csrID, serialNum uint, expired int64) (*CRL, error) {
	vs, err := daoUtils.UseTransaction(func(db *gorm.DB, serial uint, expired int64) (crl *CRL, err error) {
		// 修改 CSR 状态
		err = updateCRSStateByID(db, definition.CRSStateRevocation, csrID)
		if err != nil {
			return nil, err
		}
		// 修改证书状态
		err = updateCertificateStateByID(db, definition.CertificateStateRevocation, serial)
		if err != nil {
			return nil, err
		}
		// 插入 crl 表
		crl, err = insertNewCRL(db, serial, expired)
		if err != nil {
			return nil, err
		}
		return crl, nil
	}, []interface{}{&gorm.DB{}, serialNum, expired})
	if err != nil {
		return nil, err
	}
	crlV := vs[0].Interface().(*CRL)
	return crlV, nil
}

func GetAllCRL() ([]CRL, error) {
	crlList := make([]CRL, 0)
	err := daoUtils.GetDB().First(&crlList).Error
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("Fail to get all crls"))
		return nil, err
	}
	return crlList, nil
}
