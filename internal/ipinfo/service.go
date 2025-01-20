package ipinfo

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
)

type Service interface {
	FetchCountries() ([]models.Country, error)
	FetchCountryById(countryID string) (*models.CountryInfo, error)
	FetchCurrenciesConversionToUSD(currencyId string) (*models.CurrencyExchange, error)
	FetchCurrencies() ([]models.Currency, error)
	GetCountryByIP(ip string) (*models.IPInfo, error)
	BlockIP(ip string)
	IsBlocked(ip string) bool
}

type service struct {
	r         Repository
	blockList *BlockList
}

// NewService crea una nueva instancia del servicio.
func NewService(r Repository) Service {
	return &service{
		r:         r,
		blockList: NewBlockList(),
	}
}

// FetchCountries consulta la información de países del repositorio.
func (s *service) FetchCountries() ([]models.Country, error) {
	return s.r.FetchCountries()
}

// FetchCountryById consulta la información de un país específico del repositorio.
func (s *service) FetchCountryById(countryID string) (*models.CountryInfo, error) {
	return s.r.FetchCountryById(countryID)
}

// FetchCurrencies consulta la información de monedas del repositorio.
func (s *service) FetchCurrencies() ([]models.Currency, error) {
	return s.r.FetchCurrencies()
}

func (s *service) FetchCurrenciesConversionToUSD(currencyId string) (*models.CurrencyExchange, error) {
	return s.r.FetchCurrenciesConversionToUSD(currencyId)
}

func (s *service) GetCountryByIP(ip string) (*models.IPInfo, error) {
	// Verificacion de bloqueo de IP
	if s.blockList.IsBlocked(ip) {
		return nil, fmt.Errorf("la ip %s esta bloqueda", ip)
	}
	return s.r.GetCountryByIP(ip)
}

// BlockIP añade una IP a la lista de bloques.
func (s *service) BlockIP(ip string) {
	s.blockList.AddIP(ip)
}

// IsBlocked retorna el estado de una IP
func (s *service) IsBlocked(ip string) bool {
	return s.blockList.IsBlocked(ip)
}
