package dao

import (
	"simple_ca/src"
	"simple_ca/src/dao/utils"
	"testing"
	"time"
)

func init() {
	src.InitSetting("../../setting.json")
	utils.InitDBSetting(10, 30, time.Second*100, true)
}

func TestHasUser(t *testing.T) {
	if _, ok := HasUser("test", "test"); ok {
		t.Error("error")
	}
	if _, ok := HasUser("test", "hasPWD"); !ok {
		t.Error("error")
	}
}
