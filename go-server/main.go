package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Server представляет сервер с ресурсами
type Server struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TotalCPU     float64 `json:"total_cpu"`     // Общая мощность CPU (в ядрах)
	TotalRAM     float64 `json:"total_ram"`     // Общая оперативная память (в ГБ)
	TotalStorage float64 `json:"total_storage"` // Общее хранилище (в ГБ)
	Bandwidth    float64 `json:"bandwidth"`     // Скорость интернета (в Мбит/с)
	MonthlyPrice float64 `json:"monthly_price"` // Цена аренды за месяц (в рублях)
	PricePerHour float64 `json:"price_per_hour"`// Цена аренды за час (в рублях)
	UsedCPU      float64 `json:"used_cpu"`      // Используемая мощность CPU
	UsedRAM      float64 `json:"used_ram"`      // Используемая оперативная память
	UsedStorage  float64 `json:"used_storage"`  // Используемое хранилище
}

// LoadSnapshot представляет снимок нагрузки сервера
type LoadSnapshot struct {
	Timestamp   time.Time `json:"timestamp"`
	UsedCPU     float64   `json:"used_cpu"`
	UsedRAM     float64   `json:"used_ram"`
	UsedStorage float64   `json:"used_storage"`
}

// Lease представляет аренду ресурсов
type Lease struct {
	UserID   string  `json:"user_id"`
	ServerID string  `json:"server_id"`
	CPU      float64 `json:"cpu"`
	RAM      float64 `json:"ram"`
	Storage  float64 `json:"storage"`
}

// ServerManager управляет серверами и арендой
type ServerManager struct {
	servers     map[string]Server        // server_id -> Server
	leases      map[string]Lease         // lease_id -> Lease (lease_id = user_id:server_id)
	loadHistory map[string][]LoadSnapshot // server_id -> []LoadSnapshot
	mu          sync.Mutex               // Для потокобезопасности
}

// NewServerManager создает новый менеджер серверов
func NewServerManager() *ServerManager {
	// Предполагаем 30 дней в месяце для расчета почасовой стоимости
	const hoursInMonth = 30 * 24
	servers := map[string]Server{
		"srv1": {
			ID:           "srv1",
			Name:         "Basic Server",
			TotalCPU:     2.0,     // 2 ядра
			TotalRAM:     4.0,     // 4 ГБ
			TotalStorage: 100.0,   // 100 ГБ
			Bandwidth:    100.0,   // 100 Мбит/с
			MonthlyPrice: 5000.0,  // 5000 рублей/месяц
			PricePerHour: 5000.0 / hoursInMonth, // ~6.94 рубля/час
			UsedCPU:      0.0,
			UsedRAM:      0.0,
			UsedStorage:  0.0,
		},
		"srv2": {
			ID:           "srv2",
			Name:         "Standard Server",
			TotalCPU:     4.0,     // 4 ядра
			TotalRAM:     8.0,     // 8 ГБ
			TotalStorage: 250.0,   // 250 ГБ
			Bandwidth:    250.0,   // 250 Мбит/с
			MonthlyPrice: 10000.0, // 10000 рублей/месяц
			PricePerHour: 10000.0 / hoursInMonth, // ~13.89 рубля/час
			UsedCPU:      0.0,
			UsedRAM:      0.0,
			UsedStorage:  0.0,
		},
		"srv3": {
			ID:           "srv3",
			Name:         "Advanced Server",
			TotalCPU:     8.0,     // 8 ядер
			TotalRAM:     16.0,    // 16 ГБ
			TotalStorage: 500.0,   // 500 ГБ
			Bandwidth:    500.0,   // 500 Мбит/с
			MonthlyPrice: 20000.0, // 20000 рублей/месяц
			PricePerHour: 20000.0 / hoursInMonth, // ~27.78 рубля/час
			UsedCPU:      0.0,
			UsedRAM:      0.0,
			UsedStorage:  0.0,
		},
		"srv4": {
			ID:           "srv4",
			Name:         "Pro Server",
			TotalCPU:     12.0,    // 12 ядер
			TotalRAM:     32.0,    // 32 ГБ
			TotalStorage: 750.0,   // 750 ГБ
			Bandwidth:    750.0,   // 750 Мбит/с
			MonthlyPrice: 35000.0, // 35000 рублей/месяц
			PricePerHour: 35000.0 / hoursInMonth, // ~48.61 рубля/час
			UsedCPU:      0.0,
			UsedRAM:      0.0,
			UsedStorage:  0.0,
		},
		"srv5": {
			ID:           "srv5",
			Name:         "Elite Server",
			TotalCPU:     16.0,     // 16 ядер
			TotalRAM:     64.0,     // 64 ГБ
			TotalStorage: 1000.0,   // 1000 ГБ
			Bandwidth:    1000.0,   // 1 Гбит/с
			MonthlyPrice: 50000.0,  // 50000 рублей/месяц
			PricePerHour: 50000.0 / hoursInMonth, // ~69.44 рубля/час
			UsedCPU:      0.0,
			UsedRAM:      0.0,
			UsedStorage:  0.0,
		},
	}
	return &ServerManager{
		servers:     servers,
		leases:      make(map[string]Lease),
		loadHistory: make(map[string][]LoadSnapshot),
	}
}

// simulateLoad имитирует нагрузку на серверах с учетом арендованных ресурсов
func (sm *ServerManager) simulateLoad() {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(5 * time.Second) // Обновление каждые 5 секунд
	for range ticker.C {
		sm.mu.Lock()
		for id, server := range sm.servers {
			// Суммируем арендованные ресурсы для сервера
			leasedCPU, leasedRAM, leasedStorage := 0.0, 0.0, 0.0
			for _, lease := range sm.leases {
				if lease.ServerID == id {
					leasedCPU += lease.CPU
					leasedRAM += lease.RAM
					leasedStorage += lease.Storage
				}
			}

			// Базовая нагрузка зависит от арендованных ресурсов
			// CPU: 70-90% от арендованных + случайные флуктуации (±10%)
			// RAM: 60-80% от арендованных + случайные флуктуации (±10%)
			// Storage: 50-70% от арендованных + случайные флуктуации (±5%)
			baseCPU := leasedCPU * (0.7 + rand.Float64()*0.2)
			baseRAM := leasedRAM * (0.6 + rand.Float64()*0.2)
			baseStorage := leasedStorage * (0.5 + rand.Float64()*0.2)

			// Добавляем случайные флуктуации
			deltaCPU := baseCPU * (0.9 + rand.Float64()*0.2) // ±10%
			deltaRAM := baseRAM * (0.9 + rand.Float64()*0.2) // ±10%
			deltaStorage := baseStorage * (0.95 + rand.Float64()*0.1) // ±5%

			// Обновляем использованные ресурсы
			server.UsedCPU = deltaCPU
			server.UsedRAM = deltaRAM
			server.UsedStorage = deltaStorage

			// Ограничение, чтобы не превысить максимальные ресурсы
			if server.UsedCPU > server.TotalCPU {
				server.UsedCPU = server.TotalCPU
			}
			if server.UsedRAM > server.TotalRAM {
				server.UsedRAM = server.TotalRAM
			}
			if server.UsedStorage > server.TotalStorage {
				server.UsedStorage = server.TotalStorage
			}

			sm.servers[id] = server

			// Сохранение снимка нагрузки
			snapshot := LoadSnapshot{
				Timestamp:   time.Now(),
				UsedCPU:     server.UsedCPU,
				UsedRAM:     server.UsedRAM,
				UsedStorage: server.UsedStorage,
			}
			sm.loadHistory[id] = append(sm.loadHistory[id], snapshot)

			// Ограничение размера истории (последние 100 записей)
			if len(sm.loadHistory[id]) > 100 {
				sm.loadHistory[id] = sm.loadHistory[id][len(sm.loadHistory[id])-100:]
			}
		}
		sm.mu.Unlock()
	}
}

// getServers возвращает список всех серверов
func (sm *ServerManager) getServers(w http.ResponseWriter, r *http.Request) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	servers := make([]Server, 0, len(sm.servers))
	for _, server := range sm.servers {
		servers = append(servers, server)
	}
	json.NewEncoder(w).Encode(servers)
}

// getServerInfo возвращает информацию о конкретном сервере
func (sm *ServerManager) getServerInfo(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")
	if serverID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	server, exists := sm.servers[serverID]
	if !exists {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(server)
}

// getLoadHistory возвращает историю нагрузки для сервера
func (sm *ServerManager) getLoadHistory(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("id")
	if serverID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()
	history, exists := sm.loadHistory[serverID]
	if !exists {
		history = []LoadSnapshot{}
	}
	json.NewEncoder(w).Encode(history)
}

// leaseResources обрабатывает запрос на аренду ресурсов
func (sm *ServerManager) leaseResources(w http.ResponseWriter, r *http.Request) {
	var lease Lease
	if err := json.NewDecoder(r.Body).Decode(&lease); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if lease.ServerID == "" || lease.UserID == "" {
		http.Error(w, "UserID and ServerID are required", http.StatusBadRequest)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	server, exists := sm.servers[lease.ServerID]
	if !exists {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	// Проверка доступности ресурсов
	availableCPU := server.TotalCPU - server.UsedCPU
	availableRAM := server.TotalRAM - server.UsedRAM
	availableStorage := server.TotalStorage - server.UsedStorage

	if lease.CPU > availableCPU || lease.RAM > availableRAM || lease.Storage > availableStorage {
		http.Error(w, "Not enough resources available", http.StatusBadRequest)
		return
	}

	// Обновление использованных ресурсов
	server.UsedCPU += lease.CPU
	server.UsedRAM += lease.RAM
	server.UsedStorage += lease.Storage
	sm.servers[lease.ServerID] = server

	// Сохранение аренды
	leaseID := lease.UserID + ":" + lease.ServerID
	sm.leases[leaseID] = lease

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Resources leased successfully"})
}

// releaseResources освобождает арендованные ресурсы
func (sm *ServerManager) releaseResources(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		ServerID string `json:"server_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.ServerID == "" {
		http.Error(w, "UserID and ServerID are required", http.StatusBadRequest)
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	leaseID := req.UserID + ":" + req.ServerID
	lease, exists := sm.leases[leaseID]
	if !exists {
		http.Error(w, "No lease found for user and server", http.StatusNotFound)
		return
	}

	server, exists := sm.servers[req.ServerID]
	if !exists {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	// Освобождение ресурсов
	server.UsedCPU -= lease.CPU
	server.UsedRAM -= lease.RAM
	server.UsedStorage -= lease.Storage
	sm.servers[req.ServerID] = server

	// Удаление аренды
	delete(sm.leases, leaseID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Resources released successfully"})
}

// enableCORS добавляет заголовки CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	sm := NewServerManager()

	// Запуск имитации нагрузки
	go sm.simulateLoad()

	http.HandleFunc("/servers", enableCORS(sm.getServers))
	http.HandleFunc("/server", enableCORS(sm.getServerInfo))
	http.HandleFunc("/load-history", enableCORS(sm.getLoadHistory))
	http.HandleFunc("/lease", enableCORS(sm.leaseResources))
	http.HandleFunc("/release", enableCORS(sm.releaseResources))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}