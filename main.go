package main

import (
	"github.com/520MianXiangDuiXiang520/GinTools/email_tools"
	daoTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src"
	"time"
)

func init() {
	smtp := src.GetSetting().SMTPSetting
	email_tools.InitSMTPDialer(smtp.Host, smtp.Username, smtp.Password, smtp.Port)
	daoTools.InitDBSetting(src.GetSetting().Database, 10, 30, time.Second*100, true)
}

func main() {
	engine := gin.Default()
	defer engine.Run()
	Register(engine)
}
