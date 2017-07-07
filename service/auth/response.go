package auth

import (
	"gitlab.com/NagByte/Palette/service/common"
)

type TouchDeviceResponse struct {
	DeviceToken string `json:"deviceToken"`
	SignedIn    bool   `json:"signed_in"`
}

type WhoAmIResponse struct {
	Username string `json:"username"`
}

var (
	wrongDeviceTokenResponse        = common.ErrorJSONResponse{ErrorDescription: "wrongDeviceToken"}
	wrongUsernameOrPasswordResponse = common.ErrorJSONResponse{ErrorDescription: "wrongUsernameOrPassword"}
)
