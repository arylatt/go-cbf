package types

import "fmt"

type Product struct {
	Category   string    `json:"category"`
	StatusText string    `json:"status_text"`
	Notes      string    `json:"notes"`
	ID         string    `json:"id"`
	Allergens  Allergens `json:"allergens"`
	Dispense   string    `json:"dispense"`
	Style      string    `json:"style"`
	Name       string    `json:"name"`
	Bar        string    `json:"bar"`
	ABV        string    `json:"abv"`
}

func (p Product) Percent() string {
	if p.ABV == "" {
		return ""
	}

	return fmt.Sprintf("%s%%", p.ABV)
}
