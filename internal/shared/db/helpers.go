package db

import (
	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/lib/pq"
)

func MapPQError(err error) error {
	if err == nil {
		return nil
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err // not a pq error, return as-is
	}

	switch pqErr.Code {
	case "23505": // unique_violation
		return model.ErrUniqueViolation
	case "23503": // foreign_key_violation
		return model.ErrNoRecord
	default:
		return err // unknown error, pass through
	}
}
