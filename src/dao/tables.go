package dao

import (
	"github.com/jinzhu/gorm"
)

// 用户表
type User struct {
	gorm.Model
	Username string `gorm:"size:16;not null;unique"`
	Password string `gorm:"not null"`
	Email    string
}

// userToken 表
type UserToken struct {
	gorm.Model
	UserID     int    `gorm:"not null"`
	Token      string `gorm:"size:64;not null"`
	ExpireTime int64  `gorm:"not null"`
}

// 证书请求表
type CARequest struct {
	gorm.Model
	UserID               int    `gorm:"not null"`
	State                uint   `gorm:"not null"`
	PublicKey            string `gorm:"type:text;not null"`
	Country              string `gorm:"size:20"`
	Province             string
	Locality             string
	Organization         string
	OrganizationUnitName string
	CommonName           string
	EmailAddress         string `gorm:"not null"`
}

// 证书表
type Certificate struct {
	gorm.Model
	UserID     int   `gorm:"not null"`
	State      uint  `gorm:"not null"`
	RequestID  int   `gorm:"not null"`
	ExpireTime int64 `gorm:"not null"`
}

type CRL struct {
	gorm.Model
	CertificateID int   `gorm:"not null"`
	InputTime     int64 `gorm:"not null"`
}
