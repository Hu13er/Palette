package smsVerification

type sendVerificationRequest struct {
	PhoneNumber string `json:"phoneNumber"`
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
