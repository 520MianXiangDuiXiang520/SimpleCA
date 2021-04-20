package main

import (
	"github.com/520MianXiangDuiXiang520/GoTools/dao"
	"github.com/520MianXiangDuiXiang520/GoTools/email"
	"github.com/gin-gonic/gin"
	"simple_ca/src"
)

func init() {
	src.InitSetting("./setting.json")
	smtp := src.GetSetting().SMTPSetting
	email.InitSMTPDialer(smtp.Host, smtp.Username, smtp.Password, smtp.Port)
	_ = dao.InitDBSetting(src.GetSetting().Database)

}

func main() {
	engine := gin.Default()
	defer engine.Run()
	Register(engine)
}
