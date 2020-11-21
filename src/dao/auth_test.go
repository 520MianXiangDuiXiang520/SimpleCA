package dao

import (
	daoUtils "github.com/520MianXiangDuiXiang520/GinTools/utils/dao"
	"simple_ca/src"
	"testing"
	"time"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database, 10, 30, time.Second*100, true)
}

func TestHasUserByUP(t *testing.T) {
	if _, ok := HasUserByUP("test", "test"); ok {
		t.Error("error")
	}
	if _, ok := HasUserByUP("test", "hasPWD"); !ok {
		t.Error("error")
	}
}

func TestHasUserByID(t *testing.T) {
	u, o := HasUserByID(uint(1))
	if !o || u.ID != uint(1) {
		t.Error("Fail")
	}
	u, o = HasUserByID(uint(0))
	if o {
		t.Error("Fail")
	}
}
