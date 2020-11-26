package dao

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	utils "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/jinzhu/gorm"
)

func GetAllCerByUserID(uid uint) ([]Certificate, bool) {
	res, err := selectAllCertificateByUserID(dao_tools.GetDB(), uid)
	if err != nil {
		return nil, false
	}
	return res, true
}

func selectAllCertificateByUserID(db *gorm.DB, userID uint) ([]Certificate, error) {
	res := make([]Certificate, 0)
	err := db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		utils.ExceptionLog(err,
			fmt.Sprintf("Fail to select all Certificate by user_id: %d", userID))
		return nil, err
	}
	return res, nil
}
