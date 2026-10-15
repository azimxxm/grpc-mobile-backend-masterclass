// ============================================================================
// 04-client-streaming-telemetry: Go Backend
// Client-Streaming RPC Implementation
// ============================================================================

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

type telemetryServer struct {
	pb.UnimplementedTelemetryServiceServer
}

func (s *telemetryServer) UploadBatchMetrics(stream pb.TelemetryService_UploadBatchMetricsServer) error {
	startTime := time.Now()
	var count int32

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			// Mijoz oqimni yakunladi, umumiy hisobotni qaytaramiz:
			elapsed := time.Since(startTime).Milliseconds()
			log.Printf("Jami %d ta telemetriya voqealari qabul qilindi (%d ms)", count, elapsed)

			return stream.SendAndClose(&pb.UploadSummary{
				EventsProcessed:         count,
				Success:                 true,
				ServerProcessingTimeMs: elapsed,
			})
		}
		if err != nil {
			return err
		}

		count++
		// Voqeani tezkor tahlil qilish yoki Redis Stream-ga yozish...
		_ = event
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Port ochilmadi: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(server, &telemetryServer{})

	fmt.Println("🚀 Client-Streaming Telemetry Server listening on :50054...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Server xatosi: %v", err)
	}
}
