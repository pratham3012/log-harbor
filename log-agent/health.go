package main

import (
    "encoding/json"
    "net/http"
    "time"
	"log"
)

type HealthStatus struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Service   string    `json:"service"`
    Version   string    `json:"version"`
}

func startHealthServer(port string, serviceName string) {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        status := HealthStatus{
            Status:    "healthy",
            Timestamp: time.Now(),
            Service:   serviceName,
            Version:   "1.0.0",
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(status)
    })

    go func() {
        log.Printf("🏥 Health check server starting on port %s", port)
        if err := http.ListenAndServe(":"+port, nil); err != nil {
            log.Printf("❌ Health server error: %s", err)
        }
    }()
} 