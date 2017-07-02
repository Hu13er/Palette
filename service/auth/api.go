package auth

/*
|* API:
\************************************/

import (
	"errors"
	"io"
	"log"

	"gitlab.com/NagByte/Palette/helper"
)

var (
	// Errors:
	NoAnswerError         = io.EOF
	NotVerfiedError       = errors.New("phoneNumberNotVerified")
	WrongDeviceTokenError = errors.New("wrongDeviceToken")
	WrongUsernameOrPass   = errors.New("wrongUsernameOrPassword")
)

// TouchDevice ensures the devices node, with
// specefic UID exists in database.
// Returns deviceToken and Signedin state.
func (as *authService) TouchDevice(params map[string]interface{}) (string, bool, error) {

	query := as.db.GetQuery("touchDevice")
	maybeDeviceToken := helper.DefaultCharset.RandomStr(20)
	params["maybeDeviceToken"] = maybeDeviceToken

	result, err := as.db.QueryOne(query, params)
	// TODO: error handeling
	if err != nil && err != NoAnswerError {
		return "", false, err
	}

	deviceToken := result[0].(string)
	signedIn := result[1].(bool)

	return deviceToken, signedIn, nil
}

func (as *authService) Signup(deviceToken, username, password, verificationToken string) error {
	log.Println("In singup API.")
	defer log.Println("out singup API")
	phoneNumber, verified := as.verifier.IsVerified(verificationToken)
	if !verified {
		log.Println("verific")
		return NotVerfiedError
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
	case io.EOF:
		return WrongDeviceTokenError
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
	case io.EOF:
		return WrongUsernameOrPass
	default:
		log.Println(err)
		return err
	}
}

// Signs a device out by device Token.
func (as *authService) SignDeviceOut(deviceToken string) error {
	query := as.db.GetQuery("signDeviceOut")
	switch _, err := as.db.QueryOne(query, map[string]interface{}{"deviceToken": deviceToken}); err {
	case nil:
		return nil
	case io.EOF:
		return WrongDeviceTokenError
	default:
		log.Println(err)
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
		log.Println(err)
		return ""
	}
}
