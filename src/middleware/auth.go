package middleware

import (
	"github.com/520MianXiangDuiXiang520/ginUtils/middleware"
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
	"simple_ca/src/definition"
)

func TokenAuth(context *gin.Context) (middleware.UserBase, bool) {
	token, err := context.Cookie("SESSIONID")
	if err != nil {
		token = context.GetHeader("Token")
	}
	if len(token) != 32 {
		return nil, false
	}
	user, ok := dao.GetUserAndExtensionTime(token, definition.OneHour/2)
	if !ok {
		return nil, false
	}

	return user, true
}
