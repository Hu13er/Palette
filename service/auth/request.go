package auth

import "errors"

func checkNotEmptyStrings(strs ...string) bool {
	for _, s := range strs {
		if s == "" {
			return true
		}
	}
	return false
}

type deviceProps struct {
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Capacity string `json:"capacity"`
	OS       struct {
		Version string `json:"version"`
		Type    string `json:"type"`
	} `json:"os"`
}

func (dp *deviceProps) Validate() (map[string]interface{}, error) {
	outp := make(map[string]interface{})
	if checkNotEmptyStrings(dp.UID, dp.Name, dp.Platform, dp.Capacity, dp.OS.Type, dp.OS.Version) {
		return nil, errors.New("something is missing")
	}

	outp["uid"] = dp.UID
	outp["name"] = dp.Name
	outp["platform"] = dp.Platform
	outp["capacity"] = dp.Capacity
	outp["os_type"] = dp.OS.Type
	outp["os_version"] = dp.OS.Version

	return outp, nil
}

type signUpRequest struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	VerificationToken string `json:"verificationToken"`
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
