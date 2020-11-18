package main

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools"
	"github.com/gin-gonic/gin"
	"simple_ca/src/routes"
)

func Register(c *gin.Engine) {
	// c.Use(middleware.ApiView(), middleware2.CorsHandler())
	ginTools.URLPatterns(c, "api/ca", routes.CARegister)
	ginTools.URLPatterns(c, "api/auth", routes.AuthRegister)

}
