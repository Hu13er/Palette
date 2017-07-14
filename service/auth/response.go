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
	responseWrongDeviceToken        = common.ErrorJSONResponse{ErrorDescription: "wrongDeviceToken"}
	responseWrongUsernameOrPassword = common.ErrorJSONResponse{ErrorDescription: "wrongUsernameOrPassword"}
	responseDeviceTokenNeeded       = common.ErrorJSONResponse{ErrorDescription: "deviceTokenNeeded"}
	responseNotVerified             = common.ErrorJSONResponse{ErrorDescription: "notVerified"}
	responsePhoneNumberExists       = common.ErrorJSONResponse{ErrorDescription: "phoneNumberExists"}
	responseUsernameExists          = common.ErrorJSONResponse{ErrorDescription: "usernameExists"}
)
