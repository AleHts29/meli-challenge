package handler

import (
	"encoding/json"
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/gin-gonic/gin"
	"log"
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

		countryInfo, err := h.Service.GetCountryDataByIP(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
			if err := h.Service.BlockIP(ip); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la lista de IPs bloqueadas", "details": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "bloqueo exitoso", "count": len(req.IPs)})
	}
}

// NotifyBlockedIPs emite eventos de bloqueo a traves de Server-Sent Events (SSE)
func (h *Handler) NotifyBlockedIPs() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Configuracion de headers para SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// habilitacion de CORDS
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Canal para detectar desconexión del cliente
		disconnectChan := c.Writer.CloseNotify()

		// Suscripcion al canal de eventos
		clientChan := h.Service.SubscribeEvents()
		defer h.Service.UnsubscribeEvents(clientChan)

		// Loop para enviar eventos
		for {
			select {
			case event := <-clientChan: // Escucha eeventos en el canal principal
				_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", formatEvent(event))
				if err != nil {
					log.Printf("Error: al enviar evento: %v", err)
					return
				}
				//	Forzar la escritura del buffer al cliente
				c.Writer.Flush()
			case <-disconnectChan:
				log.Println("Cliente desconectado")
				return
			default:
			}

		}
	}
}

// formatEvent convierte un evento a JSON.
func formatEvent(event models.BlockEvent) string {
	data, _ := json.Marshal(event)
	return string(data)
}
