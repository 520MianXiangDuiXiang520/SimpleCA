package tools

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"testing"
	"time"
)

var publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAj7lVUzAXfhUs1r4L6X1+
dqZTA64Cc22Z8QRe4b2mLgaJNHx3he7Fy32a2AvkmNUD1Hz4AUyIXZvYTRM5BC+I
H86yNq4k6h3onczA+MBuzUUINt3H6diTHuoO3mpfPb9KF2WnaojURdBJ0JR4was2
2Fr2UUTgQiuw1268UjXVkc9ah6DhWGTIAlC5rrFkbY2oN0w1eQ6umcZzo+Vcs5D4
ChKosTlTGAT8k46kd9je4itQVrYSM/X9oqW1NG+HoIWmcFhKEaUTTIip+Io+o8ur
tkybWjYrJ6aL0wsghWtVLFDHwSe4cRMh/Qtvkwnbadti6Ipl/WmsyrQwG7HfXVsD
2wIDAQAB
-----END PUBLIC KEY-----`

var subject = pkix.Name{
	Country:            []string{"CN"},
	Province:           []string{"ShanXi"},
	Locality:           []string{"TaiYuan"},
	Organization:       []string{"The North University Of China"},
	OrganizationalUnit: []string{"Big Data Academy"},
	CommonName:         "SimpleCA",
}

var rootCerPath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src\\root_cer.cer"
var privatePath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src\\root_private_key.pem"
var privateKey = &rsa.PrivateKey{}
var rootCer = &x509.Certificate{}
var notBefore, notAfter = time.Now(), time.Now().Add(time.Hour * 24 * 365 * 10)
var crlPoint = []string{
	"http://39.106.168.39/crl.crl",
}

func init() {
	data, err := ioutil.ReadFile(privatePath)
	if err != nil {
		panic("can not open root private")
	}
	r, ok := DecodeRSAPrivateKey(data)
	if !ok {
		panic("can not decode root private")
	}
	rootRCer, ok := DecodePemCert(rootCerPath)
	privateKey, rootCer = r, rootRCer
}

func TestCreateIssuerRootCer(t *testing.T) {
	CreateIssuerRootCer(subject, time.Now(), time.Now().Add(time.Hour*24*365*10), privateKey, rootCerPath)
}

func TestCreateCodeSignCert(t *testing.T) {
	ok := CreateCodeSignCert(rootCer, big.NewInt(2), subject, publicKey,
		privateKey, notBefore, notAfter, crlPoint, "./test_code_sign.cer")
	if !ok {
		t.Error()
	}
}

func TestCreateSSLCert(t *testing.T) {
	ok := CreateSSLCert(rootCer, big.NewInt(3), subject, publicKey, privateKey, notBefore, notAfter, crlPoint,
		[]string{"junebao.top", "*.junebao.top"}, "./test_ssl.cer")
	if !ok {
		t.Error()
	}
}
