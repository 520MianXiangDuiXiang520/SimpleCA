package definition

const (
	// 初始化状态，用户还未提交公钥
	CRSStateInit = iota + 1
	// 审核状态，证书还未颁发
	CRSStateAuditing
	// 审核通过，证书已生成颁发
	CRSStatePass
	// 审核未通过
	CRSStateUnPass
	// 证书被撤销
	CRSStateRevocation
)

// CSR 文件内容
type CertificateSigningRequest struct {
	Country            string `json:"country"       check:"not null"`        // 国际标准组织ISO国码2位国家代号
	Province           string `json:"province"      check:"not null"`        // 省州，如 Shanxi
	Locality           string `json:"locality"      check:"not null"`        // 地区，市，如 Taiyuan
	Organization       string `json:"organization"  check:"not null"`        // 公司，组织
	CommonName         string `json:"common_name"   check:"not null"`        // FQDN(全限定域名)或姓名
	EmailAddress       string `json:"email_address" check:"not null; email"` // 邮箱
	OrganizationalUnit string `json:"organizational_unit" check:"not null"`  // 部门
}
