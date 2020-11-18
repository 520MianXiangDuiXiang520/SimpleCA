package dao

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/utils"
	"simple_ca/src"
	utils2 "simple_ca/src/dao/utils"
	"time"
)

// 根据用户名和密码检查是否有该用户存在
func HasUser(username, password string) (*utils2.User, bool) {
	user, err := selectUserByUNamePSD(username, password)
	if err != nil {
		msg := fmt.Sprintf("Fail to get user By username(%s), password(%s)", username, password)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return &user, user.ID != 0
}

func selectUserByUNamePSD(uName, pwd string) (u utils2.User, err error) {
	err = utils2.GetDB().Where("username=? AND password =?", uName, pwd).First(&u).Error
	return
}

func InsertToken(user *utils2.User, token string) (ok bool) {
	db := utils2.GetDB()
	err := db.Create(&utils2.UserToken{
		UserID:     user.ID,
		Token:      token,
		ExpireTime: time.Now().Unix() + src.GetSetting().AuthSetting.TokenExpireTime*60,
	}).Error
	if err != nil {
		msg := fmt.Sprintf("Fail to insert token; user = %v, token = %v", user, token)
		utils.ExceptionLog(err, msg)
		return false
	}
	return true
}

func InsertUser(user *utils2.User) (ok bool) {
	db := utils2.GetDB()
	err := db.Create(user).Error
	if err != nil {
		msg := fmt.Sprintf("Fail to insert new user; %v", user)
		utils.ExceptionLog(err, msg)
		return false
	}
	return true
}
