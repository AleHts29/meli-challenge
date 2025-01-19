package ipinfo

import (
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/AleHts29/meli-challenge/pkg/api"
)

type Repository interface {
	FetchCountries() ([]models.Country, error)
}

type repository struct {
	apiCountries api.Countries
}

// NewRepository crea una nueva instancia del repositorio.
func NewRepository(apiCountries api.Countries) Repository {
	return &repository{apiCountries: apiCountries}
}

// FetchCountries consulta la API de Mercado Libre para obtener información sobre países.
func (r *repository) FetchCountries() ([]models.Country, error) {
	countries, err := r.apiCountries.FetchCountries()
	if err != nil {
		return nil, err
	}
	return countries, nil
}
