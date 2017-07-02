package checkVersion

import (
	"errors"

	"gitlab.com/NagByte/Palette/db/wrapper"
)

type database struct {
	wrapper.Database
}

func (DB *database) getLatestVersion() (string, error) {
	query := DB.GetQuery("getAllVersion")
	row, err := DB.Database.QueryOne(query, nil)
	if err != nil {
		return "", err
	}
	result, ok := row[0].(string)
	if !ok {
		return "", errors.New("can not convert to string")
	}

	return result, nil
}

func (DB *database) getMinimumVersion() (string, error) {
	query := DB.GetQuery("getAllForcedVersion")
	row, err := DB.Database.QueryOne(query, nil)
	if err != nil {
		return "", err
	}
	result, ok := row[0].(string)
	if !ok {
		return "", errors.New("can not convert to string")
	}

	return result, nil
}
