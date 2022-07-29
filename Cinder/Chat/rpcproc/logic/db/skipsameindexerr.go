package db

import (
	"errors"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

func SkipSameIndexErr(err error) error {
	cmdErr := &driver.Error{}
	if errors.As(err, cmdErr) {
		if cmdErr.Code == 86 || cmdErr.Code == 85 {
			return nil
		}
	}

	return err
}
