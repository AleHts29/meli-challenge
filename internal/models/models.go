package models

type CountryInfo struct {
	Country
	DecimalSeparator        string      `json:"decimal_separator"`
	ThousandsSeparator      string      `json:"thousands_separator"`
	TimeZone                string      `json:"time_zone"`
	GeoInformation          interface{} `json:"geo_information"` // Puede ser un struct si sabes su formato
	CurrencyConversionToUSD interface{}
	States                  []State `json:"states"`
}

type Country struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	CurrencyId string `json:"currency_id"`
}

type State struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CurrencyExchange struct {
	CurrencyBase    string  `json:"currency_base"`
	CurrencyQuote   string  `json:"currency_quote"`
	Rate            float64 `json:"rate"`
	CreationDate    string  `json:"creation_date"`
	ValidUntil      string  `json:"valid_until"`
	InverseRate     float64 `json:"inv_rate"`
	LastUpdatedDate string  `json:"last_updated_date"`
}

type Currency struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `json:"decimal_places"`
}

type IPInfo struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
}
