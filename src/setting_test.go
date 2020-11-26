package src

import (
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"simple_ca/src/tools"
	"testing"
	"time"
)

func TestCreateNewCertificate(t *testing.T) {
	rootCer, rootPK := GetCARootCer()
	fmt.Println(rootPK)
	subject := pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{"Beijing"},
		Locality:           []string{"Beijing"},
		Organization:       []string{"Beijing University"},
		OrganizationalUnit: []string{"Big Data Academy"},
		CommonName:         "subject",
	}
	tools.CreateNewCertificate(&rootCer, big.NewInt(10), subject, `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAytNu1enUNQGlmYzlQYG/
r8hWoubetxf1mazDGL9SnvGjNj7F3we9lpxT8pGbYhNBh1C2SrwoEDIMy+aVKJIA
D1YxkcaRSo7H8Bri9f0zo8ZEwSY2lEw5n+dFjWuOyiD1yiCKHf074mlOMswcDYFW
edOwKVdmspw0GiRqP/9HjIl2C0xv2i6KtMgGwfKRYdEaanvFyDHxE+PdGF5m/m5+
zm1I2XS0WY2RjlIgarK/1uS9EsajFfYgG5KipiY5ZW/u7dyDzAih+LlS16cTsuwu
dj5lb2XX9x/+poka5aAW3YtG8GlVRACYv+5K9SKqUsOrifhcJxJkRSeA1FmnKRzY
sQIDAQAB
-----END PUBLIC KEY-----`, &rootPK, time.Now(), time.Now().Add(time.Hour*24*365),
		[]string{GetSetting().CRLSetting.CRLDistributionPoint}, "my.cer")
}

func TestCreateNewCRL(t *testing.T) {
	rootCer, rootPK := GetCARootCer()
	revokedCerts := []pkix.RevokedCertificate{
		{SerialNumber: big.NewInt(1), RevocationTime: time.Now()},
	}
	tools.CreateNewCRL(&rootCer, &rootPK, revokedCerts, time.Now(), time.Now().Add(time.Hour), "crl.crl")
}
