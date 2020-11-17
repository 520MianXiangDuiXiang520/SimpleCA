package dao

func CreateTable() {
	GetDB().AutoMigrate(&User{}, &UserToken{}, &CARequest{}, &Certificate{}, &CRL{})
}
