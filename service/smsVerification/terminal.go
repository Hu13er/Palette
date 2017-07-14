package smsVerification

import (
	"gitlab.com/NagByte/Palette/service/common"
)

var (
	responseInvalidPhoneNumber                 = common.ErrorJSONResponse{ErrorDescription: "invalidPhoneNumber"}
	responsePhoneNumberExists                  = common.ErrorJSONResponse{ErrorDescription: "phoneNumberExists"}
	responseWrongPhoneNumberOrVerificationCode = common.ErrorJSONResponse{ErrorDescription: "wrongPhoneNumberOrVerificationCode"}
)

type sendVerificationRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	SignUpState bool   `json:"signUpState"`
}

type verifyPhoneRequest struct {
	PhoneNumber      string `json:"phoneNumber"`
	VerificationCode string `json:"verificationCode"`
}

type verifyPhoneResponse struct {
	Token string `json:"verificationToken"`
}

type isVerifiedRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Token       string `json:"verificationToken"`
}

type isVerifiedResponse struct {
	PhoneNumber string `json:"phoneNumber"`
	Verified    bool   `json:"verified"`
}
