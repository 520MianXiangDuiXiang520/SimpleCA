package check

import (
	"testing"
)

func TestCaRequestCheck(t *testing.T) {
	if !checkPublicKey(`-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAytNu1enUNQGlmYzlQYG/
r8hWoubetxf1mazDGL9SnvGjNj7F3we9lpxT8pGbYhNBh1C2SrwoEDIMy+aVKJIA
D1YxkcaRSo7H8Bri9f0zo8ZEwSY2lEw5n+dFjWuOyiD1yiCKHf074mlOMswcDYFW
edOwKVdmspw0GiRqP/9HjIl2C0xv2i6KtMgGwfKRYdEaanvFyDHxE+PdGF5m/m5+
zm1I2XS0WY2RjlIgarK/1uS9EsajFfYgG5KipiY5ZW/u7dyDzAih+LlS16cTsuwu
dj5lb2XX9x/+poka5aAW3YtG8GlVRACYv+5K9SKqUsOrifhcJxJkRSeA1FmnKRzY
sQIDAQAB
-----END RSA PUBLIC KEY-----`) {
		t.Error("FAIL")
	}
}
