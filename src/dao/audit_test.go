package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"simple_ca/src"
	"testing"
	"time"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database)
}

func TestCreateNewCertificate(t *testing.T) {
	c, ok := CreateNewCertificate(&Certificate{
		UserID:     uint(1),
		RequestID:  uint(10),
		State:      uint(1),
		ExpireTime: time.Now().Unix(),
	})
	if !ok {
		t.Error()
	}
	fmt.Println(c)
}
