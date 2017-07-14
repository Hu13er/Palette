package smsVerification

import (
	"gitlab.com/NagByte/Palette/db/wrapper"
)

type smsVerificationDB struct {
	wrapper.Database
}

func (svd *smsVerificationDB) mergeVerificationRequest(phoneNumber, code, token string) error {

	query := svd.GetQuery("mergeVerificationRequest")
	err := svd.Exe(query, map[string]interface{}{
		"phoneNumber": phoneNumber,
		"code":        code,
		"token":       token,
	})
	return err
}

func (svd *smsVerificationDB) verifyRequest(phoneNumber, verificationCode string) (string, error) {
	query := svd.GetQuery("verifyRequest")
	result, err := svd.QueryOne(query, map[string]interface{}{
		"phoneNumber": phoneNumber,
		"code":        verificationCode,
	})
	token, _ := result[0].(string)
	return token, err
}

func (svd *smsVerificationDB) isVerified(token string) (string, bool) {
	query := svd.GetQuery("isVerified")
	result, err := svd.QueryOne(query, map[string]interface{}{
		"token": token,
	})
	if err != nil {
		return "", false
	}

	phoneNumber, _ := result[1].(string)
	ok, _ := result[0].(bool)
	return phoneNumber, ok
}
