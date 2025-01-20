package ipinfo

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/AleHts29/meli-challenge/pkg/cache"
	"log"
	"sync"
	"time"
)

const (
	BufferSizeEvents  = 256
	BufferSizeClients = 10
	CacheTime         = 5 * time.Minute
)

type Service interface {
	GetCountryDataByIP(ip string) (*models.CountryInfo, error)
	BlockIP(ip string)
	IsBlocked(ip string) bool
	SubscribeEvents() chan models.BlockEvent
	UnsubscribeEvents(clientChan chan models.BlockEvent)
}

type service struct {
	r         Repository
	blockList *BlockList
	cache     *cache.Cache
	events    chan models.BlockEvent
	mu        sync.Mutex
	clients   map[chan models.BlockEvent]struct{}
}

// NewService crea una nueva instancia del servicio.
func NewService(r Repository) Service {
	return &service{
		r:         r,
		blockList: NewBlockList(),
		cache:     cache.NewCache(CacheTime),
		events:    make(chan models.BlockEvent, BufferSizeEvents),
		clients:   make(map[chan models.BlockEvent]struct{}),
	}
}

////////////////////////////////
// *** DATA COUNTRIES ***

// GetCountryDataByIP obtiene información de un país a partir de una IP.
func (s *service) GetCountryDataByIP(ip string) (*models.CountryInfo, error) {
	// Verificar si la información ya está en la caché
	if data, ok := s.cache.Get(ip); ok {
		return data.(*models.CountryInfo), nil
	}

	// Consultar la información desde el repositorio (APIs externas)
	info, err := s.r.GetCountryByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del país para la IP: %w", err)
	}

	countryInfo, err := s.r.FetchCountryById(info.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del país: %w", err)
	}

	currencyConversion, err := s.r.FetchCurrenciesConversionToUSD(countryInfo.CurrencyId)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la cotización de la moneda: %w", err)
	}

	// Agregar la cotización al objeto de información del país
	countryInfo.CurrencyConversionToUSD = *currencyConversion

	// Guardar el resultado en la caché
	s.cache.Set(ip, countryInfo)

	return countryInfo, nil
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
