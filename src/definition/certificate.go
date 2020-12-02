package definition

const (
	// 证书使用中
	CertificateStateUsing = iota + 1

	// 证书被撤销
	CertificateStateRevocation

	// 证书过期
	CertificateStateExpired
)

const (
	CertificateTypeCodeSign = iota + 1
	CertificateTypeSSL
)
