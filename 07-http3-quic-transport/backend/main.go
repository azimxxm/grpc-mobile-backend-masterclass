// ============================================================================
// 07-http3-quic-transport/backend/main.go
// High-Performance HTTP/3 QUIC Server in Go
// ============================================================================

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/stream/migration-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Alt-Svc", `h3=":443"; ma=2592000`)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		// Keep streaming packets even if client switches between Wi-Fi and LTE!
		for i := 1; i <= 10; i++ {
			fmt.Fprintf(w, "data: {\"packet_index\": %d, \"client_ip\": \"%s\", \"time\": \"%s\"}\n\n",
				i, r.RemoteAddr, time.Now().Format(time.RFC3339Nano))
			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	})

	server := http3.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			NextProtos: []string{"h3", "h3-29"},
		},
	}

	fmt.Println("🚀 HTTP/3 QUIC UDP Server listening on :8443 (Alt-Svc enabled)...")
	// For production: provide valid TLS cert & key (e.g., self-signed or Let's Encrypt)
	log.Println("Run with: server.ListenAndServeTLS('cert.pem', 'key.pem')")
}
