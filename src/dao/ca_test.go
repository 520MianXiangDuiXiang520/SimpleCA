package dao

import (
	"fmt"
	daoUtils "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/dao_tools"
	"simple_ca/src"
	"testing"
	"time"
)

func init() {
	daoUtils.InitDBSetting(src.GetSetting().Database, 10, 30, time.Second*100, true)
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

// func TestAddPublicKeyForRequest(t *testing.T) {
// 	AddPublicKeyForRequest(1, `-----BEGIN RSA PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAytNu1enUNQGlmYzlQYG/
// r8hWoubetxf1mazDGL9SnvGjNj7F3we9lpxT8pGbYhNBh1C2SrwoEDIMy+aVKJIA
// D1YxkcaRSo7H8Bri9f0zo8ZEwSY2lEw5n+dFjWuOyiD1yiCKHf074mlOMswcDYFW
// edOwKVdmspw0GiRqP/9HjIl2C0xv2i6KtMgGwfKRYdEaanvFyDHxE+PdGF5m/m5+
// zm1I2XS0WY2RjlIgarK/1uS9EsajFfYgG5KipiY5ZW/u7dyDzAih+LlS16cTsuwu
// dj5lb2XX9x/+poka5aAW3YtG8GlVRACYv+5K9SKqUsOrifhcJxJkRSeA1FmnKRzY
// sQIDAQAB
// -----END RSA PUBLIC KEY-----`, 1)
// }
