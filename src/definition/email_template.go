package definition

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"simple_ca/src/tools"
	"strings"
)

func getEmailTemp(fail string, dict map[string]string) string {
	_, currently, _, _ := runtime.Caller(0)
	filename := path.Join(path.Dir(currently), fail)
	fmt.Println(filename)
	temp, err := ioutil.ReadFile(filename)
	if err != nil {
		tools.ExceptionLog(err, "")
		panic("Read temp Fail")
	}
	t := string(temp)
	for k, v := range dict {
		t = strings.Replace(t, "{# "+k+" #}", v, -1)
	}
	return t
}

// 证书申请成功邮件模板
func CerSuccessTemp(dict map[string]string) string {
	return getEmailTemp("./success.html", dict)
}

func CerUnPassTemp(dict map[string]string) string {
	return getEmailTemp("./fail.html", dict)
}
