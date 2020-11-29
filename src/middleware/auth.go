package middleware

import (
	"github.com/520MianXiangDuiXiang520/GinTools/gin_tools/middleware"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
)

func TokenAuth(context *gin.Context) (middleware.UserBase, bool) {
	token, err := context.Cookie("SESSIONID")
	flag := err != nil
	if !flag {
		token = context.GetHeader("Token")
	}
	if len(token) != 32 {
		return nil, false
	}
	user, ok := dao.GetUserByToken(token)

	if !ok {
		return nil, false
	}

	return user, true
}
