package models

type Country struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	CurrencyId string `json:"currency_id"`
}

type Currency struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `json:"decimal_places"`
}
