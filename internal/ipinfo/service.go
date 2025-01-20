package ipinfo

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"log"
	"sync"
)

const (
	BufferSizeEvents  = 256
	BufferSizeClients = 10
)

type Service interface {
	FetchCountries() ([]models.Country, error)
	FetchCountryById(countryID string) (*models.CountryInfo, error)
	FetchCurrenciesConversionToUSD(currencyId string) (*models.CurrencyExchange, error)
	FetchCurrencies() ([]models.Currency, error)
	GetCountryByIP(ip string) (*models.IPInfo, error)
	BlockIP(ip string)
	IsBlocked(ip string) bool
	SubscribeEvents() chan models.BlockEvent
	UnsubscribeEvents(clientChan chan models.BlockEvent)
}

type service struct {
	r         Repository
	blockList *BlockList
	events    chan models.BlockEvent
	mu        sync.Mutex
	clients   map[chan models.BlockEvent]struct{}
}

// NewService crea una nueva instancia del servicio.
func NewService(r Repository) Service {
	return &service{
		r:         r,
		blockList: NewBlockList(),
		events:    make(chan models.BlockEvent, BufferSizeEvents),
		clients:   make(map[chan models.BlockEvent]struct{}),
	}
}

////////////////////////////////
// *** DATA COUNTRIES ***

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

////////////////////////////////
// *** BLOCK_IP ***

// BlockIP añade una IP a la lista de bloqueos.
func (s *service) BlockIP(ip string) {
	s.blockList.AddIP(ip)

	// Envia notificacion a los clientes
	go func() {
		event := models.BlockEvent{IP: ip, Event: "BLOCKED"}
		s.NotifyClients(event)
		log.Printf("[INFO] Evento emitido - IP %s bloqueada", ip)
	}()
}

// IsBlocked retorna el estado de una IP
func (s *service) IsBlocked(ip string) bool {
	return s.blockList.IsBlocked(ip)
}

////////////////////////////////
// *** NOTIFICACIONES ***

// SubscribeEvents permite suscribirse al canal de eventos.
func (s *service) SubscribeEvents() chan models.BlockEvent {
	clientChan := make(chan models.BlockEvent, BufferSizeClients)
	s.mu.Lock()
	s.clients[clientChan] = struct{}{}
	s.mu.Unlock()
	return clientChan
}

// UnsubscribeEvents elimina a un cliente de la lista de suscriptores.
func (s *service) UnsubscribeEvents(clientChan chan models.BlockEvent) {
	s.mu.Lock()
	delete(s.clients, clientChan)
	close(clientChan) // cierra canal del cliente
	s.mu.Unlock()
}

// NotifyClients envía un evento a todos los clientes suscritos.
func (s *service) NotifyClients(event models.BlockEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for clientChan := range s.clients {
		clientChan <- event
	}
}
