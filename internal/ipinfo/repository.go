package ipinfo

import (
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/AleHts29/meli-challenge/pkg/api"
)

type Repository interface {
	FetchCountries() ([]models.Country, error)
	FetchCurrencies() ([]models.Currency, error)
}

type repository struct {
	apiCountries  api.Countries
	apiCurrencies api.Currencies
}

// NewRepository crea una nueva instancia del repositorio.
func NewRepository(apiCountries api.Countries, apiCurrencies api.Currencies) Repository {
	return &repository{
		apiCountries:  apiCountries,
		apiCurrencies: apiCurrencies,
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

// FetchCurrencies consulta la API de Mercado Libre para obtener información sobre las monedas.
func (r *repository) FetchCurrencies() ([]models.Currency, error) {
	currencies, err := r.apiCurrencies.FetchCurrencies()
	if err != nil {
		return nil, err
	}
	return currencies, nil
}
