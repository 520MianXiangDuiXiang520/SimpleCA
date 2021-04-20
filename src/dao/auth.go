package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"simple_ca/src"
	"simple_ca/src/tools"
	"time"
)

// 根据用户名和密码检查是否有该用户存在
func HasUserByUP(username, password string) (*User, bool) {
	user, err := selectUserByUNamePSD(username, password)
	if err != nil {
		msg := fmt.Sprintf("Fail to get user By username(%s), password(%s)", username, password)
		tools.ExceptionLog(err, msg)
		return nil, false
	}
	return &user, user.ID != 0
}

func HasUserByID(id uint) (user *User, ok bool) {
	user, err := selectUserByID(daoUtils.GetDB(), id)
	if err != nil {
		msg := fmt.Sprintf("Fail to select user by id: %d", id)
		tools.ExceptionLog(err, msg)
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

func selectUserByID(db *gorm.DB, id uint) (u *User, err error) {
	u = &User{}
	err = db.Where("id = ?", id).First(u).Error
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

func deleteTokenByUserID(db *gorm.DB, userID uint) (err error) {
	err = db.Where("user_id = ?", userID).Delete(&UserToken{}).Error
	return err
}

func DeleteTokenByUserID(uid uint) bool {
	if err := deleteTokenByUserID(daoUtils.GetDB(), uid); err != nil {
		msg := fmt.Sprintf("Fail to delete token by userID(%d)", uid)
		tools.ExceptionLog(err, msg)
		return false
	}
	return true
}

func InsertToken(user *User, token string) (ok bool) {
	// 使用事务，保证一致性
	_, err := daoUtils.UseTransaction(func(db *gorm.DB, user *User, token string) (err error) {
		err = deleteTokenByUserID(db, user.ID)
		if err != nil {
			return
		}
		return insertToken(db, user, token)
	}, []interface{}{&gorm.DB{}, user, token}, log.New(os.Stdout, "[ Transaction ] ", log.LstdFlags))

	if err != nil {
		msg := fmt.Sprintf("Fail to insert token; user = %v, token = %v", user, token)
		tools.ExceptionLog(err, msg)
		return false
	}
	return true
}

func InsertUser(user *User) (ok bool) {
	db := daoUtils.GetDB()
	err := db.Create(user).Error
	if err != nil {
		msg := fmt.Sprintf("Fail to insert new user; %v", user)
		tools.ExceptionLog(err, msg)
		return false
	}
	return true
}

func selectUserTokenByToken(db *gorm.DB, token string) (ut *UserToken, err error) {
	ut = &UserToken{}
	err = db.Where("token = ? AND expire_time > ?",
		token, time.Now().Unix()).First(ut).Error
	return ut, err
}

func GetUserByToken(token string) (user *User, ok bool) {
	ut, err := selectUserTokenByToken(daoUtils.GetDB(), token)
	if err != nil {
		msg := fmt.Sprintf("Fail to select userID by token, token: %s", token)
		tools.ExceptionLog(err, msg)
		return nil, false
	}
	id := ut.UserID
	user, err = selectUserByID(daoUtils.GetDB(), id)
	if err != nil {
		msg := fmt.Sprintf("Fail to select user by id, id: %d", id)
		tools.ExceptionLog(err, msg)
		return nil, false
	}
	return user, true
}

func GetUserByName(name string) (*User, bool) {
	u := User{}
	err := daoUtils.GetDB().Where("username = ?", name).First(&u).Error
	if err != nil {
		tools.ExceptionLog(err, fmt.Sprintf("select user by name (%s) Fail", name))
		return nil, false
	}
	return &u, true
}

func updateUserToken(db *gorm.DB, ut *UserToken) error {
	err := db.Model(&UserToken{}).Where("id = ?", ut.ID).Update(ut).Error
	if err != nil {
		tools.ExceptionLog(err, fmt.Sprintf("Fail to update userToken: %v", ut))
		return err
	}
	return nil
}

// 根据 Token 获取用户并延长过期时间
func GetUserAndExtensionTime(token string, extentTime int64) (*User, bool) {
	txFunc := func(db *gorm.DB, token string, extentTime int64) (*User, error) {
		ut, err := selectUserTokenByToken(db, token)
		if err != nil {
			return nil, err
		}

		ut.ExpireTime += extentTime
		err = updateUserToken(db, ut)
		if err != nil {
			return nil, err
		}

		user, err := selectUserByID(db, ut.UserID)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	v, err := daoUtils.UseTransaction(txFunc, []interface{}{daoUtils.GetDB(), token, extentTime}, log.New(os.Stdout, "[ Transaction ] ", log.LstdFlags))
	if err != nil {
		tools.ExceptionLog(err, fmt.Sprintf("Fail to do 'getUserAndExtensionTime'"))
		return nil, false
	}
	user := v[0].Interface().(*User)
	return user, true
}
