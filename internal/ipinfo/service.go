package ipinfo

import "github.com/AleHts29/meli-challenge/internal/models"

type Service interface {
	FetchCountries() ([]models.Country, error)
}

type service struct {
	r Repository
}

// NewService crea una nueva instancia del servicio.
func NewService(r Repository) Service {
	return &service{r: r}
}

// FetchCountries consulta la información de países del repositorio.
func (s *service) FetchCountries() ([]models.Country, error) {
	return s.r.FetchCountries()
}
