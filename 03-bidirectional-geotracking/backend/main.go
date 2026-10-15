// ============================================================================
// 03-bidirectional-geotracking: Go Backend
// Bi-directional Streaming RPC (Uber Style Radar)
// ============================================================================

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

type trackingServer struct {
	pb.UnimplementedTrackingServiceServer
}

func (s *trackingServer) LiveLocationStream(stream pb.TrackingService_LiveLocationStreamServer) error {
	for {
		// 1. Mijozdan kelayotgan GPS koordinatani o'qish
		clientLoc, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Mijoz [%s] koordinatasi: %.4f, %.4f", clientLoc.ClientId, clientLoc.Latitude, clientLoc.Longitude)

		// 2. Mijoz atrofidagi 3 ta haydovchini hisoblab, qaytarib oqizish:
		radar := &pb.NearbyDriversRadar{
			ServerTime: time.Now().Unix(),
			Drivers: []*pb.DriverInfo{
				{
					DriverId:    "driver_spark_01",
					VehicleType: "Chevrolet Spark",
					Latitude:    clientLoc.Latitude + (rand.Float64()-0.5)*0.005,
					Longitude:   clientLoc.Longitude + (rand.Float64()-0.5)*0.005,
					Heading:     45.0,
					EtaSeconds:  120,
				},
				{
					DriverId:    "driver_lacetti_02",
					VehicleType: "Chevrolet Lacetti",
					Latitude:    clientLoc.Latitude + (rand.Float64()-0.5)*0.005,
					Longitude:   clientLoc.Longitude + (rand.Float64()-0.5)*0.005,
					Heading:     180.0,
					EtaSeconds:  240,
				},
			},
		}

		if err := stream.Send(radar); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Port ochilmadi: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterTrackingServiceServer(server, &trackingServer{})

	fmt.Println("🚀 Uber-style Bi-directional Tracking Server listening on :50053...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Server xatosi: %v", err)
	}
}
