package tools

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	utils "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"os"

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
func DecryptWithDES(msgSplit []byte, key string) (r string, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			utils.ExceptionLog(errors.New("DecryptFail"),
				fmt.Sprintf("%s Decryption failed： %v", msgSplit, err))
			ok = false
		}
	}()
	keySplit := []byte(key)
	// 获取block块
	block, _ := des.NewTripleDESCipher(keySplit)
	// 创建切片
	context := make([]byte, len(msgSplit))
	// 设置解密方式
	blockMode := cipher.NewCBCDecrypter(block, keySplit[:8])
	// 解密密文到数组
	blockMode.CryptBlocks(context, msgSplit)
	// 去补码
	context, ok = PKCSUnPadding(context)
	if !ok {
		return "", false
	}
	len := int(context[0])
	return string(context[1 : len+1]), true
}

// 去码
func PKCSUnPadding(origData []byte) ([]byte, bool) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if unPadding > length {
		return nil, false
	}
	return origData[:length-unPadding], true
}

// 加密
func EncryptWithDES(msg, key string) (r []byte, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			utils.ExceptionLog(errors.New("EncryptFail"),
				fmt.Sprintf("%s Encryption failed： %v", msg, err))
			ok = false
		}
	}()
	msgSplit := []byte(msg)
	msgSplit = []byte{byte(len(msgSplit))}
	msgSplit = append(msgSplit, []byte(msg)...)
	keySplit := []byte(key)
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
	return crypt, true
}

// 补码
func PKCSPadding(origData []byte, blockSize int) []byte {
	// 计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	// 在切片后面追加char数量的byte(char)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padText...)
}

// 解码 RSA 公钥 pem 文件
func DecodeRSAPublicKey(input []byte) (interface{}, bool) {
	block, _ := pem.Decode(input)
	if block == nil || block.Type != "PUBLIC KEY" {
		utils.ExceptionLog(errors.New("DecodeRSAPublicKeyFail"),
			"failed to decode PEM block containing public key")
		return nil, false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		utils.ExceptionLog(errors.New("ParsePKIXPublicKeyFail"),
			"failed to parse PKIX public key")
		return nil, false
	}
	return pub, true
}

// 解码 RSA 私钥 pem 文件
func DecodeRSAPrivateKey(input []byte) (*rsa.PrivateKey, bool) {
	block, _ := pem.Decode(input)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		utils.ExceptionLog(errors.New("DecodeRSAPrivateKeyFail"),
			"failed to decode PEM block containing private key")
		return nil, false
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		utils.ExceptionLog(errors.New("ParsePKIXPrivateKeyFail"),
			"failed to parse PKCS1 private key")
		return nil, false
	}
	return pk, true
}

// 生成 RSA 私钥，并写入到文件
func CreateRSAPrivateKeyToFile(path string, len int) bool {
	pk, _ := rsa.GenerateKey(rand.Reader, len)
	keyOut, _ := os.Create(path)
	err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	if err != nil {
		keyOut.Close()
		utils.ExceptionLog(err, fmt.Sprintf("Fail to encode to %s", path))
		return false
	}
	keyOut.Close()
	return true
}
