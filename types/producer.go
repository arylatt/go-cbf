package types

type Producer struct {
	Location    string    `json:"location"`
	Notes       string    `json:"notes"`
	Products    []Product `json:"products"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	YearFounded string    `json:"year_founded"`
}
