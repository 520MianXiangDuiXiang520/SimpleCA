// 用于生成证书
package definition

import ()

// 签名算法
const (
	SHA256WithRSA = iota
	CCDSAWithSHA256
)

type Issuer struct {
}

type Subject struct {
}

type Validity struct {
}

type SubjectPublicKeyInfo struct {
}

type Extensions struct {
}

type Certificate struct {
	Version              uint                 // 版本
	Serial               uint                 // 序列号
	SignatureAlgorithm   uint                 // 签名算法
	Issuer               Issuer               // 颁发者
	Validity             Validity             // 有效期
	Subject              Subject              // 主体
	SubjectPublicKeyInfo SubjectPublicKeyInfo // 主体的公钥信息
	IssuerUniqueID       uint                 // 颁发者唯一 ID
	SubjectUniqueID      uint                 // 主体唯一 ID
	Extensions           Extensions           // 拓展
}

func GenerateCertificate(csr CertificateSigningRequest, pk string) {
	// subject := pkix.Name{
	//     Country:            []string{csr.Country},
	//     Province:           []string{csr.Province},
	//     Locality:           []string{csr.Locality},
	//     OrganizationalUnit: []string{csr.OrganizationalUnit},
	//     Organization:       []string{csr.Organization},
	//     CommonName:         csr.CommonName,
	// }
	// serialNumber :=
	// template := x509.Certificate{
	//
	// }

}
