package api

import (
	"errors"
	"time"

	"github.com/arylatt/go-cbf/types"
)

const (
	TimestampLayout       = "Mon Jan 02 2006 15:04:05 MST"
	TimestampLayoutLegacy = "2006-01-02 15:04:05"
)

var (
	ErrTimestampMissing = errors.New("the response does not contain a timestamp")
	ErrTimestampParse   = errors.New("failed to parse timestamp with known layouts")
)

type Response struct {
	Producers []types.Producer `json:"producers"`
	Timestamp string           `json:"timestamp"`
}

func (r *Response) Updated() (time.Time, error) {
	if r.Timestamp == "" {
		return time.Time{}, ErrTimestampMissing
	}

	if t, err := time.Parse(TimestampLayout, r.Timestamp); err == nil {
		return t, nil
	}

	if t, err := time.Parse(TimestampLayoutLegacy, r.Timestamp); err == nil {
		return t, nil
	}

	return time.Time{}, ErrTimestampParse
}
