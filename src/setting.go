package src

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/520MianXiangDuiXiang520/GoTools/json"
	path2 "github.com/520MianXiangDuiXiang520/GoTools/path"
	"io/ioutil"
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

// CRL 相关配置
type CRLSetting struct {
	CRLFileName          string `json:"crl_file_name"`          // CRL 文件名
	CRLDistributionPoint string `json:"crl_distribution_point"` // CRL 分发点
	CrlUpdateInterval    int    `json:"crl_update_interval"`    // CRL 信息更新间隔
}

// 授权访问信息相关配置
type AuthorityInfoAccess struct {
	IssuingCertificateURL string `json:"issuing_certificate_url"` // 颁发者根证书路径
}

type MySQLConn struct {
	Engine    string        `json:"engine"`
	DBName    string        `json:"db_name"`
	User      string        `json:"user"`
	Password  string        `json:"password"`
	Host      string        `json:"host"`
	Port      int           `json:"port"`
	MIdleConn int           `json:"max_idle_conn"` // 最大空闲连接数
	MOpenConn int           `json:"max_open_conn"` // 最大打开连接数
	MLifetime time.Duration `json:"max_lifetime"`  // 连接超时时间
	LogMode   bool          `json:"log_mode"`
}

type Setting struct {
	Database            *MySQLConn           `json:"database"`
	AuthSetting         *AuthSetting         `json:"auth_setting"`
	Secret              *Secret              `json:"secret"`
	SMTPSetting         *SMTPSetting         `json:"smtp_setting"`
	SiteLink            string               `json:"site_link"`
	CRLSetting          *CRLSetting          `json:"crl_setting"`
	CSRFileKey          string               `json:"csr_file_key"`
	AuthorityInfoAccess *AuthorityInfoAccess `json:"authority_info_access"`
}

var setting *Setting
var settingLock sync.Mutex
var caOnce, crlOnce sync.Once

func InitSetting(filePath string) {
	defer func() {
		if e := recover(); e != nil {
			settingLock.Unlock()
		}
	}()
	filename := filePath
	if !path2.IsAbs(filePath) {
		_, currently, _, _ := runtime.Caller(1)
		filename = path.Join(path.Dir(currently), filePath)
	}
	if setting == nil {
		settingLock.Lock()
		if setting == nil {
			err := json.FromFileLoadToObj(&setting, filename)
			if err != nil {
				panic("read setting error!")
			}
		}
		settingLock.Unlock()
	}
}

func GetSetting() *Setting {
	if setting == nil {
		panic("setting Uninitialized！")
	}
	return setting
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
		tools.ExceptionLog(err, fmt.Sprintf("open %s Fail", filename))
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
		if !tools.CreateIssuerRootCer(issuer,
			time.Now(), time.Now().Add(time.Hour*24*365*10), r, cerName) {
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
			t, ok := tools.ParseCRLUpdateTime(GetSetting().CRLSetting.CRLFileName)
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
