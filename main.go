package main

import (
	"github.com/gin-gonic/gin"
	"simple_ca/src"
	"simple_ca/src/dao/utils"
	"time"
)

func init() {
	src.InitSetting("./setting.json")
	utils.InitDBSetting(10, 30, time.Second*100, true)
}

func main() {
	engine := gin.Default()
	defer engine.Run()
	Register(engine)
}
