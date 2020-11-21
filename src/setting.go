package src

import (
	settingTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/setting_tools"
	"sync"
)

// 认证相关配置
type AuthSetting struct {
	TokenExpireTime int64 `json:"token_expire_time"` // token 过期时间，分钟
}

// 加密相关配置
type Secret struct {
	ResponseSecret string `json:"response_secret"`
}

type Setting struct {
	Database    *settingTools.DBSetting `json:"database"`
	AuthSetting *AuthSetting            `json:"auth_setting"`
	Secret      *Secret                 `json:"secret"`
}

var setting = &Setting{}
var once sync.Once

func GetSetting() Setting {
	once.Do(func() {
		settingTools.InitSetting(setting, "../setting.json")
	})
	return *setting
}
