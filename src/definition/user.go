package definition

const (
	// 管理员权限
	AdministratorRights = iota + 1
)

// 证书全量字段，用来返回给用户在线查看自己的证书
type CertificateFullAmountFields struct {
	Version               uint   `json:"version"`                 // 版本号
	SerialNumber          uint   `json:"serial_number"`           // 序列号
	Statue                uint   `json:"statue"`                  // 证书状态
	Type                  uint   `json:"type"`                    // 证书类型
	NotBefore             int64  `json:"not_before"`              // 证书起始时间
	NotAfter              int64  `json:"not_after"`               // 证书过期时间
	Subject               string `json:"subject"`                 // 证书主体
	Issuer                string `json:"issuer"`                  // 证书颁发者
	SignatureAlgorithm    string `json:"signature_algorithm"`     // 签名算法
	PublicKeyAlgorithm    string `json:"public_key_algorithm"`    // 公钥算法
	CRLDistributionPoints string `json:"crl_distribution_points"` // CRL 分发点
	KeyUsage              string `json:"key_usage"`               // 密钥用法
	ExtKeyUsage           string `json:"ext_key_usage"`           // 增强型密钥用法
	PublicKey             string `json:"public_key"`              // 公钥1
	DownloadLink          string `json:"download_link"`           // 证书下载路径
	DNSName               string `json:"dns_name"`                // 使用者可选名称
}
