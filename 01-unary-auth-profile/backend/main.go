// ============================================================================
// 01-unary-auth-profile: Go Backend
// REST (net/http + JSON) vs gRPC (google.golang.org/grpc) Taqqoslashi
// ============================================================================

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

// ----------------------------------------------------------------------------
// 1. gRPC Server Implementation
// ----------------------------------------------------------------------------

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Email va parol kiritilishi shart")
	}

	return &pb.LoginResponse{
		AccessToken:  "grpc_jwt_token_example_123456",
		RefreshToken: "grpc_refresh_token_7890",
		ExpiresIn:    3600,
	}, nil
}

func (s *authServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.UserProfile, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.NotFound, "Foydalanuvchi topilmadi")
	}

	// 10x ixcham ikkilik formatda uzatiladi:
	return &pb.UserProfile{
		UserId:     req.UserId,
		FullName:   "Azimjon Pulatov",
		Email:      "azimjon@example.com",
		AvatarUrl:  "https://cdn.example.com/avatars/user.jpg",
		IsVerified: true,
		CreatedAt:  time.Now().Unix(),
	}, nil
}

// ----------------------------------------------------------------------------
// 2. An'anaviy REST Handler (Taqqoslash uchun)
// ----------------------------------------------------------------------------

func restProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	profile := map[string]interface{}{
		"user_id":     "usr_123",
		"full_name":   "Azimjon Pulatov",
		"email":       "azimjon@example.com",
		"avatar_url":  "https://cdn.example.com/avatars/user.jpg",
		"is_verified": true,
		"created_at":  time.Now().Unix(),
	}
	_ = json.NewEncoder(w).Encode(profile)
}

func main() {
	// gRPC Serverni ishga tushirish (Port :50051)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("gRPC port ochilmadi: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &authServer{})

	go func() {
		fmt.Println("🚀 gRPC Server listening on :50051 (Binary Protobuf over HTTP/2)...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server to'xtadi: %v", err)
		}
	}()

	// REST Serverni ishga tushirish (Port :8080)
	http.HandleFunc("/v1/users/profile", restProfileHandler)
	fmt.Println("🌐 REST Server listening on :8080 (Text JSON over HTTP/1.1)...")
	_ = http.ListenAndServe(":8080", nil)
}
