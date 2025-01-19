package handler

import (
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler define el manejador HTTP para las solicitudes relacionadas con IPs y países.
type Handler struct {
	Service ipinfo.Service
}

// NewHandler crea un nuevo manejador para las solicitudes relacionadas con IPs y países.
func NewHandler(s ipinfo.Service) *Handler {
	return &Handler{Service: s}
}

// GetCountryByIP devuelve información sobre un país a partir de una IP.
func (h *Handler) GetCountryByIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Param("ip")
		if ip == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar una IP"})
			return
		}

		// Obtiene el countryId de un pais mediante una IP
		info, err := h.Service.GetCountryByIP(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener información del país para la IP: " + err.Error()})
			return
		}

		// Obtiene información asociada al país.
		countryInfo, err := h.Service.FetchCountryById(info.CountryCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener información detallada del país: " + err.Error()})
			return
		}

		// Obtiene la contizacion en dolares.
		currencyConversion, err := h.Service.FetchCurrenciesConversionToUSD(countryInfo.CurrencyId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la cotización de la moneda en dólares: " + err.Error()})
			return
		}

		// Se agrega la cotización de la moneda al objeto `countryInfo`.
		countryInfo.CurrencyConversionToUSD = *currencyConversion

		// Se devuelve la información del país en formato JSON.
		c.JSON(http.StatusOK, countryInfo)
	}
}

// GetCountries devuelve una lista de países.
func (h *Handler) GetCountries() gin.HandlerFunc {
	return func(c *gin.Context) {
		countries, err := h.Service.FetchCountries()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, countries)
	}
}

// GetCurrency devuelve la moneda asociada a una un Pais
func (h *Handler) GetCurrency() gin.HandlerFunc {
	return func(c *gin.Context) {
		currency, err := h.Service.FetchCurrencies()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, currency)
	}
}
