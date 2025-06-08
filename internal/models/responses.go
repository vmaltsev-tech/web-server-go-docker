package models

// HealthResponse представляет ответ health check endpoint
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// InfoResponse представляет ответ info endpoint
type InfoResponse struct {
	Message     string `json:"message"`
	Environment string `json:"environment"`
	Port        string `json:"port"`
}

// MetricsResponse представляет ответ metrics endpoint
type MetricsResponse struct {
	RequestCount int    `json:"request_count"`
	Uptime       string `json:"uptime"`
	StartTime    string `json:"start_time"`
}
