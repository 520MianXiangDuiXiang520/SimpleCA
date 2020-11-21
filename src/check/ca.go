package check

import (
	ginTools "github.com/520MianXiangDuiXiang520/GinTools/gin_tools"
	"github.com/520MianXiangDuiXiang520/GinTools/log_tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"simple_ca/src/message"
)

func checkPublicKey(str string) bool {
	re, err := regexp.Compile("^-----BEGIN [A-Z ]*PUBLIC KEY-----[a-zA-Z\\n0-9/+]*-----END [A-Z ]*PUBLIC KEY-----$")
	if err != nil {
		utils.ExceptionLog(err, "compile public key Fail")
		return false
	}
	return re.MatchString(str)
}

func CaRequestCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	_ = req.(*message.CaCodeSignatureRequestReq)
	// if !checkPublicKey(r.PublicKey) {
	//     return ginTools.ParamErrorRespHeader, errors.New("")
	// }
	return http.StatusOK, nil
}

func CaCsrCheck(ctx *gin.Context, req ginTools.BaseReqInter) (ginTools.BaseRespInter, error) {
	return http.StatusOK, nil
}
