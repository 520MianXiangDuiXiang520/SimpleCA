package utils

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"github.com/jinzhu/gorm"
	"simple_ca/src"
	"simple_ca/src/dao"
	"testing"
	"time"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database, 10, 30, time.Second*100, true)
}

func TestUseTransaction(t *testing.T) {
	r, _ := daoUtils.UseTransaction(func(db *gorm.DB, id uint, token string) error {
		return db.Create(&dao.UserToken{
			UserID:     id,
			Token:      token,
			ExpireTime: 99,
		}).Error
	}, []interface{}{&gorm.DB{}, uint(1), "19999"})
	for _, v := range r {
		fmt.Println(v.Interface())
	}
}
