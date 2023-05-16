package types

import (
	"encoding/json"
	"errors"
	"strings"
)

var ErrAllergenUnmarshal = errors.New("could not unmarshal allergen information")

type Allergen bool

type Allergens struct {
	Gluten      Allergen `json:"gluten,omitempty"`
	Sulphites   Allergen `json:"sulphites,omitempty"`
	Egg         Allergen `json:"egg,omitempty"`
	Nuts        Allergen `json:"nuts,omitempty"`
	Mustard     Allergen `json:"mustard,omitempty"`
	Celery      Allergen `json:"celery,omitempty"`
	Peanuts     Allergen `json:"peanuts,omitempty"`
	Fish        Allergen `json:"fish,omitempty"`
	Sesame      Allergen `json:"sesame,omitempty"`
	Milk        Allergen `json:"milk,omitempty"`
	Lupins      Allergen `json:"lupins,omitempty"`
	Soybeans    Allergen `json:"soybeans,omitempty"`
	Molluscs    Allergen `json:"molluscs,omitempty"`
	Crustaceans Allergen `json:"crustaceans,omitempty"`
}

func (a *Allergen) UnmarshalJSON(b []byte) (err error) {
	for _, f := range []func(b []byte) (err error){a.unmarshalInt, a.unmarshalString, a.unmarshalBool} {
		if err := f(b); err == nil {
			return nil
		}
	}

	return ErrAllergenUnmarshal
}

func (a *Allergen) unmarshalString(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}

	if strings.TrimSpace(s) == "" {
		*a = false
		return
	}

	*a = true
	return
}

func (a *Allergen) unmarshalInt(b []byte) (err error) {
	var i int
	if err = json.Unmarshal(b, &i); err != nil {
		return
	}

	if i > 0 {
		*a = true
		return
	}

	*a = false
	return
}

func (a *Allergen) unmarshalBool(b []byte) (err error) {
	var b1 bool
	if err = json.Unmarshal(b, &b1); err != nil {
		return
	}

	*a = Allergen(b1)
	return
}
