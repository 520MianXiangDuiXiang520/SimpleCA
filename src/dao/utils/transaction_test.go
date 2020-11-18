package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"simple_ca/src"
	"testing"
	"time"
)

func init() {
	src.InitSetting("../../../setting.json")
	InitDBSetting(10, 30, time.Second*100, true)
}

func TestUseTransaction(t *testing.T) {
	r, _ := UseTransaction(func(db *gorm.DB, id uint, token string) error {
		return db.Create(&UserToken{
			UserID:     id,
			Token:      token,
			ExpireTime: 99,
		}).Error
	}, []interface{}{&gorm.DB{}, 1, "19999"})
	for _, v := range r {
		fmt.Println(v.Interface())
	}
}
