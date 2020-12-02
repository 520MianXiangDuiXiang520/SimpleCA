package tools

import (
	"fmt"
	"io/ioutil"

	// "simple_ca/src"
	"testing"
)

func TestEncryptBySHA256(t *testing.T) {
	r := "simpleCA"
	s := HashBySHA256([]string{r})
	fmt.Println(s)
}

func TestHashByMD5(t *testing.T) {
	r := "simpleCA"
	s := HashBySHA256([]string{r})
	fmt.Println(s)
}

// func TestEncryptWithDES(t *testing.T) {
// 	msg, _ := EncryptWithDES("123", src.GetSetting().Secret.ResponseSecret)
// 	fmt.Println(msg)
// 	m, _ := DecryptWithDES(msg, src.GetSetting().Secret.ResponseSecret)
// 	fmt.Println(m)
// 	if m != "123" {
// 		t.Error("FAIL")
// 	}
// }

func TestCreateRSAPrivateKeyToFile(t *testing.T) {
	if !CreateRSAPrivateKeyToFile("TestCreateRSAPrivateKeyToFile.pem", 2048) {
		t.Error()
	}
}

func TestDecodeRSAPrivateKey(t *testing.T) {
	data, _ := ioutil.ReadFile("TestCreateRSAPrivateKeyToFile.pem")
	pk, ok := DecodeRSAPrivateKey(data)
	if !ok {
		t.Error()
	}
	fmt.Println(pk)
}
