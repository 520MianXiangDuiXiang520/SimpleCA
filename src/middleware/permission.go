package middleware

import (
	"github.com/gin-gonic/gin"
	"simple_ca/src/dao"
	"simple_ca/src/definition"
)

func AdminPermit(ctx *gin.Context) bool {
	user, ok := ctx.Get("user")
	if !ok {
		return false
	}
	u := user.(*dao.User)
	return u.Authority == definition.AdministratorRights
}
