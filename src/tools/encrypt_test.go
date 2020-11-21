package tools

import (
	"fmt"
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

func TestEncryptWithDES(t *testing.T) {
	msg := EncryptWithDES("123")
	fmt.Println(msg)
	m := DecryptWithDES(msg)
	fmt.Println(m)
	if m != "123" {
		t.Error("FAIL")
	}
}
