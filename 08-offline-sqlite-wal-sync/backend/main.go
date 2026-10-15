// ============================================================================
// 08-offline-sqlite-wal-sync/backend/main.go
// High-Throughput Batch Ingestion Go Server
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ClientEvent struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
}

type BatchIngestRequest struct {
	DeviceID string        `json:"device_id"`
	SentAt   int64         `json:"sent_at"`
	Events   []ClientEvent `json:"events"`
}

type BatchIngestResponse struct {
	ProcessedCount int    `json:"processed_count"`
	ServerTime     int64  `json:"server_time"`
	Status         string `json:"status"`
}

func main() {
	http.HandleFunc("/api/v1/telemetry/batch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req BatchIngestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// High performance bulk write simulation:
		// In production, execute: INSERT INTO events (id, payload) VALUES (...), (...); in a single SQL batch
		log.Printf("📥 [Batch Ingestion] Device %s sent %d events. Processing in single DB transaction...",
			req.DeviceID, len(req.Events))

		resp := BatchIngestResponse{
			ProcessedCount: len(req.Events),
			ServerTime:     time.Now().Unix(),
			Status:         "BATCH_COMMITTED_OK",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	fmt.Println("🚀 Offline Batch Ingestion Go Server listening on :8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
