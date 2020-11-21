package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/jinzhu/gorm"
	"simple_ca/src"
	"time"
)

// 根据用户名和密码检查是否有该用户存在
func HasUserByUP(username, password string) (*User, bool) {
	user, err := selectUserByUNamePSD(username, password)
	if err != nil {
		msg := fmt.Sprintf("Fail to get user By username(%s), password(%s)", username, password)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return &user, user.ID != 0
}

func HasUserByID(id uint) (user *User, ok bool) {
	user, err := selectUserByID(id)
	if err != nil {
		msg := fmt.Sprintf("Fail to select user by id: %d", id)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	if user.ID == 0 {
		return nil, false
	}
	return user, true
}

func selectUserByUNamePSD(uName, pwd string) (u User, err error) {
	err = daoUtils.GetDB().Where("username=? AND password =?", uName, pwd).First(&u).Error
	return
}

func selectUserByID(id uint) (u *User, err error) {
	u = &User{}
	err = daoUtils.GetDB().Where("id = ?", id).First(u).Error
	return u, err
}

func insertToken(db *gorm.DB, user *User, token string) (err error) {
	err = db.Create(&UserToken{
		UserID:     user.ID,
		Token:      token,
		ExpireTime: time.Now().Unix() + src.GetSetting().AuthSetting.TokenExpireTime*60,
	}).Error
	return err
}

func deleteTokenByUser(db *gorm.DB, user *User) (err error) {
	err = db.Where("user_id = ?", user.ID).Delete(&UserToken{}).Error
	return err
}

func InsertToken(user *User, token string) (ok bool) {
	// 使用事务，保证一致性
	_, err := daoUtils.UseTransaction(func(db *gorm.DB, user *User, token string) (err error) {
		err = deleteTokenByUser(db, user)
		if err != nil {
			return
		}
		return insertToken(db, user, token)
	}, []interface{}{&gorm.DB{}, user, token})

	if err != nil {
		msg := fmt.Sprintf("Fail to insert token; user = %v, token = %v", user, token)
		utils.ExceptionLog(err, msg)
		return false
	}
	return true
}

func InsertUser(user *User) (ok bool) {
	db := daoUtils.GetDB()
	err := db.Create(user).Error
	if err != nil {
		msg := fmt.Sprintf("Fail to insert new user; %v", user)
		utils.ExceptionLog(err, msg)
		return false
	}
	return true
}

func selectUserIDByToken(token string) (uid uint, err error) {
	ut := &UserToken{}
	err = daoUtils.GetDB().Where("token = ? AND expire_time > ?",
		token, time.Now().Unix()).First(ut).Error
	uid = ut.UserID
	return uid, err
}

func GetUserByToken(token string) (user *User, ok bool) {
	id, err := selectUserIDByToken(token)
	if err != nil {
		msg := fmt.Sprintf("Fail to select userID by token, token: %s", token)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	user, err = selectUserByID(id)
	if err != nil {
		msg := fmt.Sprintf("Fail to select user by id, id: %d", id)
		utils.ExceptionLog(err, msg)
		return nil, false
	}
	return user, true
}
