package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"simple_ca/src"
	"simple_ca/src/definition"
	"testing"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database)
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

func TestGetUserByName(t *testing.T) {
	u, ok := GetUserByName("ggdjs")
	fmt.Println(u, ok)
}

func TestGetUserAndExtensionTime(t *testing.T) {
	u, e := GetUserAndExtensionTime("848ae7a63594b5f7cdd00d8ccad30e75", definition.OneHour/2)
	if !e {
		t.Error()
	}
	fmt.Println(u)
}
