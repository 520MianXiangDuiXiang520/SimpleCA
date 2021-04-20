package dao

import (
	"fmt"
	daotools "github.com/520MianXiangDuiXiang520/GoTools/dao"
	"github.com/jinzhu/gorm"
	"simple_ca/src"
	"simple_ca/src/definition"
	"simple_ca/src/tools"
)

func GetCertificateFullAmountFieldsUser(user *User) ([]definition.CertificateFullAmountFields, bool) {
	issuerSetting := src.GetSetting().Secret.CAIssuerInfo
	issuer := fmt.Sprintf("%s, %s, %s, %s, %s, %s", issuerSetting.CommonName,
		issuerSetting.OrganizationalUnit, issuerSetting.Organization,
		issuerSetting.Locality, issuerSetting.Province, issuerSetting.Country)
	cers, err := selectAllCertificateByUserID(daotools.GetDB(), user.ID)
	if err != nil {
		return nil, false
	}
	result := make([]definition.CertificateFullAmountFields, len(cers))
	for i, v := range cers {
		csr, ok := GetCRSByID(v.RequestID)
		if !ok {
			return nil, false
		}
		subject := fmt.Sprintf("%s, %s, %s, %s, %s, %s", csr.CommonName,
			csr.OrganizationUnitName, csr.Organization, csr.Locality,
			csr.Province, csr.Country)
		cName := tools.GetCertificateFileName(v.ID, user.ID, user.Username)
		result[i] = definition.CertificateFullAmountFields{
			Issuer:                issuer,
			Subject:               subject,
			SerialNumber:          v.ID,
			Statue:                v.State,
			NotAfter:              v.ExpireTime,
			NotBefore:             v.CreatedAt.Unix(),
			SignatureAlgorithm:    "SHA256WithRSA",
			PublicKeyAlgorithm:    "RSA",
			KeyUsage:              "KeyUsageDigitalSignature | KeyUsageCertSign",
			ExtKeyUsage:           "ExtKeyUsageAny",
			CRLDistributionPoints: src.GetSetting().CRLSetting.CRLDistributionPoint,
			PublicKey:             csr.PublicKey,
			DownloadLink:          src.GetSetting().Secret.DownloadLink + cName,
			Type:                  csr.Type,
			DNSName:               csr.DnsNames,
		}
	}
	return result, true
}

func selectAllCertificateByUserID(db *gorm.DB, userID uint) ([]Certificate, error) {
	res := make([]Certificate, 0)
	err := db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		tools.ExceptionLog(err,
			fmt.Sprintf("Fail to select all Certificate by user_id: %d", userID))
		return nil, err
	}
	return res, nil
}
