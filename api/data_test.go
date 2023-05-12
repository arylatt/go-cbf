package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseUpdated(t *testing.T) {
	tests := []struct {
		timestamp    string
		expectedTime time.Time
		expectedErr  error
	}{
		{
			timestamp:    "Wed May 10 2023 22:08:23 BST",
			expectedTime: time.Date(2023, time.May, 10, 22, 8, 23, 0, time.Local),
			expectedErr:  nil,
		},
		{
			timestamp:    "2023-05-10 22:08:23",
			expectedTime: time.Date(2023, time.May, 10, 22, 8, 23, 0, time.UTC),
			expectedErr:  nil,
		},
		{
			timestamp:    "",
			expectedTime: time.Time{},
			expectedErr:  ErrTimestampMissing,
		},
		{
			timestamp:    "a",
			expectedTime: time.Time{},
			expectedErr:  ErrTimestampParse,
		},
	}

	for _, test := range tests {
		time, err := (&Response{Timestamp: test.timestamp}).Updated()

		assert.ErrorIs(t, err, test.expectedErr)
		assert.True(t, test.expectedTime.Equal(time), "%s should match %s for %s", time.String(), test.expectedTime.String(), test.timestamp)
	}
}
