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
