package utils

func CreateTable() {
	GetDB().AutoMigrate(&User{}, &UserToken{}, &CARequest{}, &Certificate{}, &CRL{})
}
