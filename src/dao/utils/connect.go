package utils

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GinTools/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"simple_ca/src"
	"strings"
	"time"
)

type DBConnector struct {
	setting *src.DBSetting
}

func (conn *DBConnector) NewConnect() *gorm.DB {
	connURI := ""
	switch strings.ToLower(conn.setting.Engine) {
	case "mysql":
		connURI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
			conn.setting.User, conn.setting.Password,
			conn.setting.Host, conn.setting.Port, conn.setting.DBName)
	case "":
		panic("engine is nil")
	default:
		panic("unrecognized database engine")
	}

	db, err := gorm.Open(conn.setting.Engine, connURI)
	if err != nil {
		msg := fmt.Sprintf("Fail to open db, URI is %s", connURI)
		utils.ExceptionLog(err, msg)
		panic(err)
	}
	return db
}

// 设置连接池参数
func setup(maxIdle, maxOpen int, maxLifeTime time.Duration, logMode bool) {
	db.DB().SetMaxIdleConns(maxIdle)        // 最大空闲连接数
	db.DB().SetMaxOpenConns(maxOpen)        // 最大连接数
	db.DB().SetConnMaxLifetime(maxLifeTime) // 设置连接空闲超时
	db.LogMode(logMode)
}

var (
	dbConnector = DBConnector{}
	db          = &gorm.DB{}
)

func InitDBSetting(maxIdle, maxOpen int, maxLifeTime time.Duration, logMode bool) {
	dbConnector.setting = src.GetSetting().Database
	db = dbConnector.NewConnect()
	setup(maxIdle, maxOpen, maxLifeTime, logMode)
}

func GetDB() *gorm.DB {
	if dbConnector.setting == nil || dbConnector.setting.Engine == "" {
		panic("Database configuration is not loaded")
	}
	if db == nil {
		panic("Public db is nil")
	}
	if err := db.DB().Ping(); err != nil {
		_ = db.Close()
		db = dbConnector.NewConnect()
	}
	return db
}
