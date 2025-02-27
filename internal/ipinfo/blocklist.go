package ipinfo

import "sync"

// BlockList maneja la lista de IPs bloqueadas.
type BlockList struct {
	blockedIPs map[string]bool
	mu         sync.Mutex
}

// NewBlockList crea una nueva instancia de la lista de IPs bloqueadas.
func NewBlockList() *BlockList {
	return &BlockList{
		blockedIPs: make(map[string]bool),
	}
}

// GetAll retorna la lista de IPs bloqueadas
func (bl *BlockList) GetAll() []string {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	blockedIPs := make([]string, 0, len(bl.blockedIPs))
	for ip := range bl.blockedIPs {
		blockedIPs = append(blockedIPs, ip)
	}
	return blockedIPs
}

// AddIP bloquea una IP
func (bl *BlockList) AddIP(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	bl.blockedIPs[ip] = true
}

// IsBlocked verifica si una IP está bloqueada.
func (bl *BlockList) IsBlocked(ip string) bool {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	return bl.blockedIPs[ip]
}

// RemoveIP desbloquea una IP
func (bl *BlockList) RemoveIP(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	delete(bl.blockedIPs, ip)
}
