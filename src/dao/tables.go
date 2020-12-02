package dao

import (
	"github.com/jinzhu/gorm"
)

// 用户表
type User struct {
	gorm.Model
	Username  string `gorm:"size:16;not null;unique"`
	Password  string `gorm:"not null"`
	Email     string
	Authority uint // 权限：1 管理员权限
}

func (u *User) GetID() int {
	return int(u.ID)
}

// userToken 表
type UserToken struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	Token      string `gorm:"size:64;not null"`
	ExpireTime int64  `gorm:"not null"`
}

// 证书请求表
type CARequest struct {
	gorm.Model
	UserID               uint `gorm:"not null"`
	State                uint `gorm:"not null"`
	Type                 uint
	PublicKey            string `gorm:"type:text;not null"`
	Country              string `gorm:"size:20"`
	Province             string
	Locality             string
	Organization         string
	OrganizationUnitName string
	CommonName           string
	EmailAddress         string `gorm:"not null"`
	DnsNames             string
}

// 证书表
type Certificate struct {
	gorm.Model
	UserID     uint  `gorm:"not null"`
	State      uint  `gorm:"not null"`
	RequestID  uint  `gorm:"not null"`
	ExpireTime int64 `gorm:"not null"`
}

type CRL struct {
	gorm.Model
	CertificateID uint  `gorm:"not null"`
	InputTime     int64 `gorm:"not null"`
}
