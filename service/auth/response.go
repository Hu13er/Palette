package auth

type TouchDeviceResponse struct {
	DeviceToken string `json:"deviceToken"`
	SignedIn    bool   `json:"signed_in"`
}

type WhoAmIResponse struct {
	Username string `json:"username"`
}
