package api

import (
	"encoding/json"
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"net/http"
)

type Countries interface {
	FetchCountries() ([]models.Country, error)
	CountryById(countryID string) (*models.Country, error)
}

type apiData struct {
	accessToken string
	apiUrl      string
}

// NewCountries crea una nueva instancia de la estructura Contries.
func NewCountries(accessToken, apiUrl string) Countries {
	return &apiData{
		accessToken: accessToken,
		apiUrl:      apiUrl,
	}
}

// FetchCountries consulta la API de Mercado Libre para obtener información sobre países.
func (a *apiData) FetchCountries() ([]models.Country, error) {
	url := fmt.Sprintf("%s/classified_locations/countries", a.apiUrl)

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
	var countries []models.Country
	if err = json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return countries, nil
}

// CountryById consulta la API de Mercado Libre para obtener información sobre un país específico.
func (a *apiData) CountryById(countryID string) (*models.Country, error) {
	// Implementar la lógica para consultar la API de Mercado Libre
	// y devolver la información del país específico
	return nil, nil
}
