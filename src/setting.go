package src

import (
	"encoding/json"
	"github.com/520MianXiangDuiXiang520/GinTools/utils"
	"os"
	"path"
	"runtime"
)

type DBSetting struct {
	Engine   string `json:"engine"`
	DBName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type AuthSetting struct {
	TokenExpireTime int64 `json:"token_expire_time"` // token 过期时间，分钟
}

type Setting struct {
	Database    *DBSetting   `json:"database"`
	AuthSetting *AuthSetting `json:"auth_setting"`
}

func (setting *Setting) load(path string) {
	fp, err := os.Open(path)
	if err != nil {
		utils.ExceptionLog(err, "Fail to open setting")
		panic(err)
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&setting)
	if err != nil {
		utils.ExceptionLog(err, "Fail to decode json setting")
		panic(err)
	}
}

var setting *Setting

func InitSetting(fPath string) {
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), fPath)
	s := Setting{}
	s.load(filename)
	setting = &s
}

func GetSetting() Setting {
	if setting == nil {
		panic("Configuration file not loaded")
	}
	return *setting
}
