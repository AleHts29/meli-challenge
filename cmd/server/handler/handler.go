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
