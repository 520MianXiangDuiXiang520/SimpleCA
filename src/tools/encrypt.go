package tools

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"simple_ca/src"
	"strings"
)

// 使用 SHA256 哈希
func HashBySHA256(sList []string) (res string) {
	s := strings.Join(sList, "")
	hash := sha256.New()
	hash.Write([]byte(s))
	hash.Sum(nil)
	bytes := hash.Sum(nil)
	res = hex.EncodeToString(bytes)
	return
}

// 使用 MD5 做哈希
func HashByMD5(strList []string) (h string) {
	r := strings.Join(strList, "")
	hash := md5.New()
	hash.Write([]byte(r))
	return hex.EncodeToString(hash.Sum(nil))
}

// 解密
func DecryptWithDES(msg string) string {
	keySplit := []byte(src.GetSetting().Secret.ResponseSecret)
	msgSplit, _ := base64.StdEncoding.DecodeString(msg)
	fmt.Println(msgSplit, msg)
	// 获取block块
	block, _ := des.NewTripleDESCipher(keySplit)
	// 创建切片
	context := make([]byte, len(msgSplit))
	// 设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, keySplit[:8])
	// 解密密文到数组
	blockMode.CryptBlocks(context, msgSplit)
	// 去补码
	context = PKCSUnPadding(context)
	len := int(context[0])
	return string(context[1 : len+1])
}

// 去码
func PKCSUnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:length-unPadding]
}

// 加密
func EncryptWithDES(msg string) string {
	msgSplit := []byte(msg)
	msgSplit = []byte{byte(len(msgSplit))}
	msgSplit = append(msgSplit, []byte(msg)...)
	keySplit := []byte(src.GetSetting().Secret.ResponseSecret)
	// 获取block块
	block, _ := des.NewTripleDESCipher(keySplit)
	// 补码
	msgSplit = PKCSPadding(msgSplit, block.BlockSize())
	// 设置加密方式为 3DES  使用3条56位的密钥对数据进行三次加密
	blockMode := cipher.NewCBCEncrypter(block, keySplit[:8])
	// 创建明文长度的数组
	crypt := make([]byte, len(msgSplit))
	// 加密明文
	blockMode.CryptBlocks(crypt, msgSplit)
	return base64.StdEncoding.EncodeToString(crypt)
}

// 补码
func PKCSPadding(origData []byte, blockSize int) []byte {
	// 计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	// 在切片后面追加char数量的byte(char)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}
