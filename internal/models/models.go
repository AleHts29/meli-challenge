package models

type Country struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	CurrencyId string `json:"currency_id"`
}
