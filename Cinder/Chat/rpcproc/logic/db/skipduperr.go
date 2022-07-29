package db

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// SkipDupErr 返回原错误，只是忽略Dup错误而变成nil
func SkipDupErr(err error) error {
	if err == nil {
		return nil
	}

	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code != 11000 {
				return err
			}
		}
		return nil
	}

	var be mongo.BulkWriteException
	if errors.As(err, &be) {
		for _, we := range be.WriteErrors {
			if we.Code != 11000 {
				return err
			}
		}
		return nil
	}

	return err
}
