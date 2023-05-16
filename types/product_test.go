package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductPercent(t *testing.T) {
	tests := []struct {
		abv      string
		expected string
	}{
		{
			"",
			"",
		},
		{
			"69",
			"69%",
		},
		{
			"42.0",
			"42.0%",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, (Product{ABV: test.abv}).Percent())
	}
}
