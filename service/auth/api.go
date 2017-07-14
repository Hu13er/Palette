package auth

/*
|* API:
\************************************/

import (
	"errors"
	"io"

	"gitlab.com/NagByte/Palette/helper"
)

var (
	// Errors:
	ErrNotFound            = io.EOF
	ErrNotVerfied          = errors.New("phoneNumberNotVerified")
	ErrWrongDeviceToken    = errors.New("wrongDeviceToken")
	ErrWrongUsernameOrPass = errors.New("wrongUsernameOrPassword")
	ErrUsernameExists      = errors.New("usernameExists")
	ErrPhoneNumberExists   = errors.New("phoneNumberExists")
)

// TouchDevice ensures the devices node, with
// specefic UID exists in database.
// Returns deviceToken and Signedin state.
func (as *authService) TouchDevice(params map[string]interface{}) (string, bool, string, error) {

	query := as.db.GetQuery("touchDevice")
	maybeDeviceToken := helper.DefaultCharset.RandomStr(20)
	params["maybeDeviceToken"] = maybeDeviceToken

	result, err := as.db.QueryOne(query, params)
	// TODO: error handeling
	if err != nil && err != ErrNotFound {
		return "", false, "", err
	}

	deviceToken := result[0].(string)
	signedIn, _ := result[1].(bool)
	username, _ := result[2].(string)

	return deviceToken, signedIn, username, nil
}

func (as *authService) Signup(deviceToken, username, password, verificationToken string) error {
	phoneNumber, verified := as.verifier.IsVerified(verificationToken)
	if !verified {
		return ErrNotVerfied
	}

	if !as.IsUniqueUsername(username) {
		return ErrUsernameExists
	}

	if !as.IsUniquePhoneNumber(phoneNumber) {
		return ErrPhoneNumberExists
	}

	query := as.db.GetQuery("signUp")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{
		"deviceToken": deviceToken,
		"username":    username,
		"password":    password,
		"phoneNumber": phoneNumber,
	}); err {
	case nil:
		return nil
	default:
		return err
	}
}

// Signs a device in by deviceToken.
func (as *authService) SignDeviceIn(deviceToken, username, password string) error {

	query := as.db.GetQuery("signDeviceIn")
	params := map[string]interface{}{
		"deviceToken": deviceToken,
		"username":    username,
		"password":    password,
	}

	switch _, err := as.db.QueryOne(query, params); err {
	case nil:
		return nil
	case ErrNotFound:
		return ErrWrongUsernameOrPass
	default:
		return err
	}
}

// Signs a device out by device Token.
func (as *authService) SignDeviceOut(deviceToken string) error {
	query := as.db.GetQuery("signDeviceOut")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{"deviceToken": deviceToken}); err {
	case nil:
		return nil
	default:
		return err
	}
}

// Gets deviceToken and returns Username
// Null string if not logged in.
func (as *authService) WhoAmI(deviceToken string) string {
	query := as.db.GetQuery("whoAmI")
	switch result, err := as.db.QueryOne(query, map[string]interface{}{"deviceToken": deviceToken}); err {
	case nil:
		return result[0].(string)
	default:
		return ""
	}
}

func (as *authService) ensureDeviceToken(deviceToken string) bool {
	query := as.db.GetQuery("ensureDeviceToken")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{"deviceToken": deviceToken}); err {
	case nil:
		return true
	default:
		return false
	}
}

func (as *authService) IsUniquePhoneNumber(phoneNumber string) bool {
	query := as.db.GetQuery("isUniquePhoneNumber")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{"phoneNumber": phoneNumber}); err {
	case nil:
		return false
	default:
		return true
	}
}

func (as *authService) IsUniqueUsername(username string) bool {
	query := as.db.GetQuery("isUniqueUsername")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{"username": username}); err {
	case nil:
		return false
	default:
		return true
	}
}
