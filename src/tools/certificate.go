package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"runtime"
	"time"
)

func DecodePemCert(p string) (*x509.Certificate, bool) {
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), p)
	pemTmp, err := ioutil.ReadFile(filename)
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("read %s Fail", filename))
		return nil, false
	}
	certBlock, _ := pem.Decode(pemTmp)
	if certBlock == nil {
		ExceptionLog(err, "pem decode fail")
		return nil, false
	}
	// 证书解析
	certBody, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		ExceptionLog(err, "fail to parse cert")
		return nil, false
	}
	return certBody, true
}

func createNewCertificate(rootCer, template *x509.Certificate,
	publicKey string, pk *rsa.PrivateKey, p string) bool {

	var c []byte
	var err error
	// 生成根证书
	if rootCer == nil {
		c, err = x509.CreateCertificate(rand.Reader, template, template, &pk.PublicKey, pk)
	} else {
		pub, ok := DecodeRSAPublicKey([]byte(publicKey))
		if !ok {
			return false
		}
		c, err = x509.CreateCertificate(rand.Reader, template, rootCer, pub, pk)
	}
	if err != nil {
		ExceptionLog(err, "Failed to create certificate")
		return false
	}
	certOut, err := os.Create(p)
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Failed to create %s", p))
		return false
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: c})
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Failed to encode pem"))
		return false
	}
	certOut.Close()
	return true
}

// 生成颁发者根证书
func CreateIssuerRootCer(issuer pkix.Name, notBefore,
	notAfter time.Time, pk *rsa.PrivateKey, p string) bool {
	template := &x509.Certificate{
		Version:            1,
		SerialNumber:       big.NewInt(1),
		Subject:            issuer,
		Issuer:             issuer,
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		NotBefore:          notBefore,
		NotAfter:           notAfter,
		IsCA:               true,
		// IssuingCertificateURL:
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	return createNewCertificate(nil, template, "", pk, p)
}

// 生成代码签名证书
func CreateCodeSignCert(rootCer *x509.Certificate, serialN *big.Int, subject pkix.Name,
	publicKey string, pk *rsa.PrivateKey, notBefore, notAfter time.Time,
	CRLDistributionPoint []string, p string) bool {
	template := &x509.Certificate{
		Version:               1,
		SerialNumber:          serialN,
		Subject:               subject,
		Issuer:                subject,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		CRLDistributionPoints: CRLDistributionPoint,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
	}
	return createNewCertificate(rootCer, template, publicKey, pk, p)
}

// 生成 SSL 证书
func CreateSSLCert(rootCer *x509.Certificate, serialN *big.Int, subject pkix.Name,
	publicKey string, pk *rsa.PrivateKey, notBefore, notAfter time.Time,
	CRLDistributionPoint []string, DNSNames []string, p string) bool {
	template := &x509.Certificate{
		Version:               1,
		SerialNumber:          serialN,
		Subject:               subject,
		Issuer:                subject,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		CRLDistributionPoints: CRLDistributionPoint,
		DNSNames:              DNSNames,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	return createNewCertificate(rootCer, template, publicKey, pk, p)
}

// 生成 CRL
func CreateNewCRL(cer *x509.Certificate, pk *rsa.PrivateKey,
	revokedCerts []pkix.RevokedCertificate, now, expiry time.Time, fileName string) bool {
	crlBytes, err := cer.CreateCRL(rand.Reader, pk, revokedCerts, now, expiry)
	if err != nil {
		ExceptionLog(err, "Fail to create CRL")
		return false
	}
	certOut, err := os.Create(fileName)
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Failed to create %s", fileName))
		return false
	}
	err = pem.Encode(certOut, &pem.Block{Type: "X509 CRL", Bytes: crlBytes})
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Failed to encode pem"))
		return false
	}
	certOut.Close()
	return true
}

// 从 CRL 文件读取上次更新时间
func ParseCRLUpdateTime(filePath string) (int64, bool) {
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), filePath)
	crlF, err := ioutil.ReadFile(filename)
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Fail to read csr: %s", filename))
		return 0, false
	}
	crl, err := x509.ParseCRL(crlF)
	if err != nil {
		ExceptionLog(err, fmt.Sprintf("Fail to parse csr: %s", filename))
		return 0, false
	}
	return crl.TBSCertList.ThisUpdate.Unix(), true
}

func GetCertificateFileName(serial, userID uint, userName string) string {
	return HashByMD5([]string{
		fmt.Sprintf("%d, %s", serial, userName),
		fmt.Sprintf("%d, %s", userID, userName),
		fmt.Sprintf("%d. %d. %s", userID, serial, userName),
	}) + ".cer"
}
