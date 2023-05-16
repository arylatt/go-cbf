package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllergenUnmarshalJSON(t *testing.T) {
	tests := []struct {
		json     []byte
		err      error
		allergen Allergen
	}{
		{
			[]byte("{}"),
			ErrAllergenUnmarshal,
			Allergen(false),
		},
		{
			[]byte("0"),
			nil,
			Allergen(false),
		},
		{
			[]byte("1"),
			nil,
			Allergen(true),
		},
		{
			[]byte("\"\""),
			nil,
			Allergen(false),
		},
		{
			[]byte("\"1\""),
			nil,
			Allergen(true),
		},
		{
			[]byte("false"),
			nil,
			Allergen(false),
		},
		{
			[]byte("true"),
			nil,
			Allergen(true),
		},
	}

	for _, test := range tests {
		a := toAllergenPtr(Allergen(false))

		assert.ErrorIs(t, test.err, a.UnmarshalJSON(test.json))
		assert.Equal(t, test.allergen, *a)
	}
}

func toAllergenPtr(a Allergen) *Allergen {
	return &a
}
