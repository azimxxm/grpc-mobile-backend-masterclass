// ============================================================================
// 02-server-streaming-chat: Go Backend
// Server-Streaming RPC Implementation
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *chatServer) SubscribeChannel(req *pb.ChannelSubscriptionRequest, stream pb.ChatService_SubscribeChannelServer) error {
	log.Printf("Mijoz kanalga ulandi: %s (oxirgi id: %d)", req.ChannelId, req.LastMessageId)

	// Simulyatsiya: Har 1.5 soniyada yangi xabar oqizish
	for i := 1; i <= 20; i++ {
		// Mijoz aloqani uzganligini tekshirish (Cancellation)
		if stream.Context().Err() != nil {
			log.Printf("Mijoz kanalni yopdi, oqim to'xtatildi")
			return stream.Context().Err()
		}

		msg := &pb.ChatMessage{
			MessageId:  int64(i),
			ChannelId:  req.ChannelId,
			SenderName: fmt.Sprintf("Foydalanuvchi #%d", (i%5)+1),
			Text:       fmt.Sprintf("Salom! Bu jonli oqimdagi xabar #%d", i),
			Timestamp:  time.Now().Unix(),
			IsPinned:   i == 1,
		}

		if err := stream.Send(msg); err != nil {
			return err
		}

		time.Sleep(1500 * time.Millisecond)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Port ochilmadi: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterChatServiceServer(server, &chatServer{})

	fmt.Println("🚀 Chat gRPC Server-Streaming listening on :50052...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Server xatosi: %v", err)
	}
}
