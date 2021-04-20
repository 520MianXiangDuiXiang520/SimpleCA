package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"simple_ca/src"
	"testing"
	"time"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database)
}

func TestCreateNewCRS(t *testing.T) {
	newCSR := &CARequest{
		UserID:               uint(1),
		State:                uint(1),
		PublicKey:            "",
		Country:              "CN",
		Province:             "Shanxi",
		Locality:             "Taiyuan",
		Organization:         "NUC",
		OrganizationUnitName: "Big Data Academy",
		CommonName:           "blog.junebao.top",
		EmailAddress:         "15364968962@163.com",
	}
	csr, ok := CreateNewCRS(newCSR)
	if !ok {
		t.Error("Fail")
	}
	if csr.ID == 0 {
		t.Error("result is nil")
	}
	fmt.Println(csr)
}

func TestHasCRSByID(t *testing.T) {
	if HasCRSByID(0) {
		t.Error("FAIL")
	}
	if !HasCRSByID(1) {
		t.Error("FAIL")
	}
}

func TestCreateNewCRL(t *testing.T) {
	crl, err := CreateNewCRL(uint(1), uint(9), time.Now().Unix())
	if err != nil {
		t.Error()
	}
	fmt.Println(crl)
}

func TestGetAllCRL(t *testing.T) {
	fmt.Println(GetAllCRL())
}
