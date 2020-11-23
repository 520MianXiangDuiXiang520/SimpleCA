package src

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	settingTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools/setting_tools"
	utils "github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"runtime"
	"simple_ca/src/definition"
	"simple_ca/src/tools"
	"sync"
	"time"
)

// 认证相关配置
type AuthSetting struct {
	TokenExpireTime int64 `json:"token_expire_time"` // token 过期时间，分钟
}

// 加密和证书相关配置
type Secret struct {
	ResponseSecret           string                                `json:"response_secret"`
	CARootPrivateKeyName     string                                `json:"ca_root_private_key_name"`
	CARootPrivateKeyLen      int                                   `json:"ca_root_private_key_len"`
	CARootCerName            string                                `json:"ca_root_cer_name"`
	CAIssuerInfo             *definition.CertificateSigningRequest `json:"ca_issuer_info"`
	CertificateEffectiveTime int64                                 `json:"certificate_effective_time"` // 证书有效时长，单位天
}

type Setting struct {
	Database    *settingTools.DBSetting `json:"database"`
	AuthSetting *AuthSetting            `json:"auth_setting"`
	Secret      *Secret                 `json:"secret"`
}

var setting = &Setting{}
var once, caOnce sync.Once

func GetSetting() Setting {
	once.Do(func() {
		settingTools.InitSetting(setting, "../setting.json")
	})
	return *setting
}

var CARootCer = &x509.Certificate{}
var CARootPrivateKey = &rsa.PrivateKey{}

// 加载 CA 私钥和根证书
func loadCAKey() (rootRCer *x509.Certificate, rootRPK *rsa.PrivateKey) {
	// 获取私钥
	pkName := GetSetting().Secret.CARootPrivateKeyName
	_, err := os.Stat(pkName)
	if !tools.HasThisFile(pkName) {
		// 私钥不存在，创建
		if !tools.CreateRSAPrivateKeyToFile(pkName, GetSetting().Secret.CARootPrivateKeyLen) {
			panic("CAPrivateKeyAcquisitionFailed")
		}
	}
	_, currently, _, _ := runtime.Caller(1)
	filename := path.Join(path.Dir(currently), pkName)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		utils.ExceptionLog(err, fmt.Sprintf("open %s Fail", filename))
		panic("CAPrivateKeyAcquisitionFailed")
	}
	r, ok := tools.DecodeRSAPrivateKey(data)
	rootRPK = r
	if !ok {
		panic("CAPrivateKeyAcquisitionFailed")
	}

	// 加载证书
	s := GetSetting().Secret.CAIssuerInfo
	issuer := pkix.Name{
		Country:            []string{s.Country},
		Province:           []string{s.Province},
		Locality:           []string{s.Locality},
		Organization:       []string{s.Organization},
		OrganizationalUnit: []string{s.OrganizationalUnit},
		CommonName:         s.CommonName,
	}
	cerName := GetSetting().Secret.CARootCerName
	if !tools.HasThisFile(cerName) {
		// 新建证书
		if !tools.CreateNewCertificate(nil, big.NewInt(1), issuer,
			"", r, time.Now(), time.Now().Add(time.Hour*24*365*10), cerName) {
			panic("FailedToCreateRootCertificate")
		}
	}
	// 读根证书
	rootRCer, ok = tools.DecodePemCert(cerName)
	if !ok {
		panic("DecodeRootCertificateFail")
	}
	return
}

func GetCARootCer() (x509.Certificate, rsa.PrivateKey) {
	caOnce.Do(func() {
		CARootCer, CARootPrivateKey = loadCAKey()
	})
	return *CARootCer, *CARootPrivateKey
}
