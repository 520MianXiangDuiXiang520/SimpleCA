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
	"sync/atomic"
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
	UserCerPath              string                                `json:"user_cer_path"`
	CARootCerName            string                                `json:"ca_root_cer_name"`
	CAIssuerInfo             *definition.CertificateSigningRequest `json:"ca_issuer_info"`
	CertificateEffectiveTime int64                                 `json:"certificate_effective_time"` // 证书有效时长，单位天
	DownloadLink             string                                `json:"download_link"`              // 证书下载路径
}

// SMTP 连接相关配置
type SMTPSetting struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CRLSetting struct {
	CRLFileName          string `json:"crl_file_name"`          // CRL 文件名
	CRLDistributionPoint string `json:"crl_distribution_point"` // CRL 分发点
	CrlUpdateInterval    int    `json:"crl_update_interval"`    // CRL 信息更新间隔
}

type Setting struct {
	Database    *settingTools.DBSetting `json:"database"`
	AuthSetting *AuthSetting            `json:"auth_setting"`
	Secret      *Secret                 `json:"secret"`
	SMTPSetting *SMTPSetting            `json:"smtp_setting"`
	SiteLink    string                  `json:"site_link"`
	CRLSetting  *CRLSetting             `json:"crl_setting"`
}

var setting = &Setting{}
var once, caOnce, crlOnce sync.Once

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
			"", r, time.Now(), time.Now().Add(time.Hour*24*365*10),
			[]string{GetSetting().CRLSetting.CRLDistributionPoint}, cerName) {
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

var crlUpdateTimeNextTime int64

// 获取下一次更新 CRL 的时间

func GetNextUpdateCRLTime() int64 {
	if atomic.LoadInt64(&crlUpdateTimeNextTime) == 0 {
		crlOnce.Do(func() {
			t, ok := tools.ParseCRLUpdateTime("../crl.crl")
			if !ok {
				atomic.StoreInt64(&crlUpdateTimeNextTime, time.Now().Unix())
			} else {
				atomic.StoreInt64(&crlUpdateTimeNextTime, t)
			}
		})
	}
	return atomic.LoadInt64(&crlUpdateTimeNextTime)
}

func SetNextUpdateCRLTime(n int64) {
	atomic.StoreInt64(&crlUpdateTimeNextTime, n)
}
