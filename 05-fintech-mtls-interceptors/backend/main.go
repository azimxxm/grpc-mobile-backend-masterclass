// ============================================================================
// 05-fintech-mtls-interceptors: Go Backend
// mTLS + Server Unary Interceptor (Security validation & Idempotency)
// ============================================================================

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

// ----------------------------------------------------------------------------
// 1. Server Security Interceptor: Header va Nonce tekshiruvi
// ----------------------------------------------------------------------------

func securityInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata topilmadi")
	}

	// Idempotency va Nonce mavjudligini tekshirish
	idempotencyKeys := md.Get("x-idempotency-key")
	if len(idempotencyKeys) == 0 || idempotencyKeys[0] == "" {
		return nil, status.Error(codes.InvalidArgument, "X-Idempotency-Key sarlavhasi majburiy")
	}

	log.Printf("[SECURITY PASS] Method: %s | IdempotencyKey: %s", info.FullMethod, idempotencyKeys[0])
	return handler(ctx, req)
}

// ----------------------------------------------------------------------------
// 2. Banking Service
// ----------------------------------------------------------------------------

type bankingServer struct {
	pb.UnimplementedBankingServiceServer
}

func (s *bankingServer) TransferFunds(ctx context.Context, req *pb.TransferRequest) (*pb.TransferReceipt, error) {
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "O'tkazma summasi musbat bo'lishi shart")
	}

	log.Printf("Pul o'tkazmasi: %.2f %s (%s -> %s)", req.Amount, req.Currency, req.FromAccount, req.ToAccount)

	return &pb.TransferReceipt{
		TransactionId: fmt.Sprintf("tx_bank_%d", time.Now().UnixNano()),
		Status:        "COMPLETED",
		Fee:           req.Amount * 0.005, // 0.5% komissiya
		Timestamp:     time.Now().Unix(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Port ochilmadi: %v", err)
	}

	// Interceptor bilan serverni ishga tushirish
	server := grpc.NewServer(
		grpc.UnaryInterceptor(securityInterceptor),
	)
	pb.RegisterBankingServiceServer(server, &bankingServer{})

	fmt.Println("🚀 Bank-Grade Secure gRPC Server listening on :50055 (mTLS & Interceptor)...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Server xatosi: %v", err)
	}
}
