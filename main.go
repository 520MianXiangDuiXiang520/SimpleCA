package main

import (
	daoTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"github.com/gin-gonic/gin"
	"simple_ca/src"
	"time"
)

func init() {
	daoTools.InitDBSetting(src.GetSetting().Database, 10, 30, time.Second*100, true)
}

func main() {
	engine := gin.Default()
	defer engine.Run()
	Register(engine)
}
