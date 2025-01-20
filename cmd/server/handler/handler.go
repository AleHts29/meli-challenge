package handler

import (
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/gin-gonic/gin"
	"net"
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

		// Validacion de formato de la IP
		if net.ParseIP(ip) == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La IP proporcionada no es válida", "ip": ip})
			return
		}

		// Verifica si la IP esta bloqueada.
		if h.Service.IsBlocked(ip) {
			c.JSON(http.StatusForbidden, gin.H{"error": "IP está bloqueda, no es posible visualizar la informacion"})
			return
		}

		// TODO: Validar solo los paises en los que opera MELI

		// TODO: Implementar un mecanismo de cache para mejorar la performance en consultas repetidas

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

// BlockIPs bloquea una o varias IPs para evitar que se consulte informacion del pais de origen
func (h *Handler) BlockIPs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			IPs []string `json:"ip" binding:"required"`
		}

		// Intentar parsear el cuerpo de la solicitud
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Debe proporcionar una IP válida en el cuerpo de la solicitud", "details": err.Error()})
			return
		}

		if len(req.IPs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Lista de IPs vacia, agregar las IPs que desea bloquear"})
			return
		}

		// Validar formato de la lista de IPs
		for _, ip := range req.IPs {
			if net.ParseIP(ip) == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "La IP proporcionada no es válida", "ip": ip})
				return
			}
		}

		// Se bloquea la IP
		for _, ip := range req.IPs {
			h.Service.BlockIP(ip)
		}

		c.JSON(http.StatusOK, gin.H{"message": "bloqueo exitoso", "count": len(req.IPs)})
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
