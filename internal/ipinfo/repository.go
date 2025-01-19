package ipinfo

import (
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/AleHts29/meli-challenge/pkg/api"
	"github.com/AleHts29/meli-challenge/pkg/store"
)

type Repository interface {
	FetchCountries() ([]models.Country, error)
	FetchCountryById(countryID string) (*models.CountryInfo, error)
	FetchCurrenciesConversionToUSD(currencyId string) (*models.CurrencyExchange, error)
	FetchCurrencies() ([]models.Currency, error)
	GetCountryByIP(ip string) (*models.IPInfo, error)
}

type repository struct {
	apiCountries  api.Countries
	apiCurrencies api.Currencies
	ipStore       store.IpStore
}

// NewRepository crea una nueva instancia del repositorio.
func NewRepository(apiCountries api.Countries, apiCurrencies api.Currencies, ipStore store.IpStore) Repository {
	return &repository{
		apiCountries:  apiCountries,
		apiCurrencies: apiCurrencies,
		ipStore:       ipStore,
	}
}

// FetchCountries consulta la API de Mercado Libre para obtener información sobre países.
func (r *repository) FetchCountries() ([]models.Country, error) {
	countries, err := r.apiCountries.FetchCountries()
	if err != nil {
		return nil, err
	}
	return countries, nil
}

// FetchCountryById consulta la API de Mercado Libre para obtener información sobre un país específico.
func (r *repository) FetchCountryById(countryID string) (*models.CountryInfo, error) {
	country, err := r.apiCountries.FetchCountryById(countryID)
	if err != nil {
		return nil, err
	}
	return country, nil
}

// FetchCurrencies consulta la API de Mercado Libre para obtener información sobre las monedas.
func (r *repository) FetchCurrencies() ([]models.Currency, error) {
	currencies, err := r.apiCurrencies.FetchCurrencies()
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

func (r *repository) FetchCurrenciesConversionToUSD(currencyId string) (*models.CurrencyExchange, error) {
	currency, err := r.apiCurrencies.FetchCurrenciesConversionToUSD(currencyId)
	if err != nil {
		return nil, err
	}
	return currency, nil
}

func (r *repository) GetCountryByIP(ip string) (*models.IPInfo, error) {
	byIP, err := r.ipStore.GetCountryByIP(ip)
	if err != nil {
		return nil, err
	}
	return byIP, nil
}
