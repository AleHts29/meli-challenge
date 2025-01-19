package api

import (
	"encoding/json"
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"net/http"
)

type Currencies interface {
	FetchCurrencies() ([]models.Currency, error)
	CurrenciesById(countryID string) (*models.Currency, error)
}

type apiCurrencies struct {
	//accessToken string
	apiUrl string
}

// NewCountries crea una nueva instancia de la estructura Contries.
func NewCurrencies(accessToken, apiUrl string) Currencies {
	return &apiCurrencies{
		//accessToken: accessToken,
		apiUrl: apiUrl,
	}
}

// FetchCountries consulta la API de Mercado Libre para obtener información sobre países.
func (a *apiCurrencies) FetchCurrencies() ([]models.Currency, error) {
	url := fmt.Sprintf("%s/currencies", a.apiUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	//req.Header.Set("Authorization", "Bearer "+a.accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching countries: status code %d", resp.StatusCode)
	}
	var currencies []models.Currency
	if err = json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return currencies, nil
}

// CountryById consulta la API de Mercado Libre para obtener información sobre un país específico.
func (a *apiCurrencies) CurrenciesById(countryID string) (*models.Currency, error) {

	// Implementar la lógica para consultar la API de Mercado Libre
	// y devolver la información del país específico

	return nil, nil
}
