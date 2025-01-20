package api

import (
	"encoding/json"
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"net/http"
)

type Countries interface {
	FetchCountries() ([]models.Country, error)
	FetchCountryById(countryID string) (*models.CountryInfo, error)
}

type apiCountries struct {
	//accessToken string
	apiUrl string
}

// NewCountries crea una nueva instancia de la estructura Contries.
func NewCountries(accessToken, apiUrl string) Countries {
	return &apiCountries{
		//accessToken: accessToken,
		apiUrl: apiUrl,
	}
}

// FetchCountries consulta la API de Mercado Libre para obtener información sobre países.
func (a *apiCountries) FetchCountries() ([]models.Country, error) {
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

// FetchCountryById consulta la API de Mercado Libre para obtener información sobre un país específico.
func (a *apiCountries) FetchCountryById(countryID string) (*models.CountryInfo, error) {
	url := fmt.Sprintf("%s/classified_locations/countries/%s", a.apiUrl, countryID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching countryByID: status code %d", resp.StatusCode)
	}

	var country models.CountryInfo
	if err = json.NewDecoder(resp.Body).Decode(&country); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &country, nil
}
