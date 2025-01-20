package ipinfo

import (
	"encoding/json"
	"fmt"
	"github.com/AleHts29/meli-challenge/internal/models"
	"github.com/AleHts29/meli-challenge/pkg/cache"
	"log"
	"os"
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
	BlockIP(ip string) error
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
	filePath  string // ruta del archivo para guardar estados de la aplicacion
}

// NewService crea una nueva instancia del servicio.
func NewService(r Repository, filePath string) Service {
	service := &service{
		r:         r,
		blockList: NewBlockList(),
		cache:     cache.NewCache(CacheTime),
		events:    make(chan models.BlockEvent, BufferSizeEvents),
		clients:   make(map[chan models.BlockEvent]struct{}),
		filePath:  filePath,
	}
	go service.loadBlockedIPs()
	return service
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

	// Obtener la lista de países en los que opera MELI
	countries, err := s.r.FetchCountries()
	if err != nil {
		return nil, fmt.Errorf("error al obtener la lista de países: %w", err)
	}

	// Validar si el CountryCode está en la lista de países
	if !s.isValidCountry(info.CountryCode, countries) {
		return nil, fmt.Errorf("el país con código '%s' no es válido o no opera MercadoLibre", info.CountryCode)
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

// isValidCountry verifica si el countryCode está en la lista de países.
func (s *service) isValidCountry(countryCode string, countries []models.Country) bool {
	for _, country := range countries {
		if country.ID == countryCode {
			return true
		}
	}
	return false
}

////////////////////////////////
// *** BLOCK_IP ***

// BlockIP añade una IP a la lista de bloqueos.
func (s *service) BlockIP(ip string) error {
	s.blockList.AddIP(ip)

	// Guardar estado de la aplicacion en el archivo
	err := s.saveBlockedIPs()
	if err != nil {
		return err
	}

	// Envia notificacion de bloqueo a clientes
	go func() {
		event := models.BlockEvent{IP: ip, Event: "BLOCKED"}
		s.notifyClients(event)
		log.Printf("[INFO] Evento emitido - IP %s bloqueada", ip)
	}()

	return nil
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

// notifyClients envía un evento a todos los clientes suscritos.
func (s *service) notifyClients(event models.BlockEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for clientChan := range s.clients {
		clientChan <- event
	}
}

////////////////////////////////
// *** APP_STATE ***

// SaveBlockedIPs guarda la lista de IPs bloqueadas en un archivo.
func (s *service) saveBlockedIPs() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(s.blockList.GetAll())
	if err != nil {
		return err
	}
	_, err = file.Write(data)

	return err
}

// LoadBlockedIPs carga la lista de IPs bloqueadas desde un archivo.
func (s *service) loadBlockedIPs() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// si el archivo no existe inicializa lista vacia
			s.blockList = NewBlockList()
			return nil
		}
		return err
	}
	defer file.Close()

	var blockedIPs []string
	if err := json.NewDecoder(file).Decode(&blockedIPs); err != nil {
		return err
	}

	for _, ip := range blockedIPs {
		s.blockList.AddIP(ip)
	}

	return nil
}
