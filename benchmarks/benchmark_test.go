// ============================================================================
// benchmarks/benchmark_test.go
// REST JSON vs Protobuf Serialization & Deserialization Benchmark
// ============================================================================

package benchmarks

import (
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	pb "github.com/azimxxm/grpc-mobile-backend-masterclass/gen/go/masterclass/v1"
)

type JSONProfile struct {
	UserID     string `json:"user_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	AvatarURL  string `json:"avatar_url"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  int64  `json:"created_at"`
}

var sampleJSON = JSONProfile{
	UserID:     "usr_987654321",
	FullName:   "Azimjon Pulatov",
	Email:      "azimjon@example.com",
	AvatarURL:  "https://cdn.example.com/avatars/user_avatar_highres_image.jpg",
	IsVerified: true,
	CreatedAt:  time.Now().Unix(),
}

var samplePB = pb.UserProfile{
	UserId:     "usr_987654321",
	FullName:   "Azimjon Pulatov",
	Email:      "azimjon@example.com",
	AvatarUrl:  "https://cdn.example.com/avatars/user_avatar_highres_image.jpg",
	IsVerified: true,
	CreatedAt:  time.Now().Unix(),
}

// 1. JSON Serialization (Encoding)
func BenchmarkJSONMarshal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&sampleJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 2. Protobuf Serialization (Encoding)
func BenchmarkProtobufMarshal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(&samplePB)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 3. JSON Deserialization (Decoding)
func BenchmarkJSONUnmarshal(b *testing.B) {
	data, _ := json.Marshal(&sampleJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p JSONProfile
		if err := json.Unmarshal(data, &p); err != nil {
			b.Fatal(err)
		}
	}
}

// 4. Protobuf Deserialization (Decoding)
func BenchmarkProtobufUnmarshal(b *testing.B) {
	data, _ := proto.Marshal(&samplePB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p pb.UserProfile
		if err := proto.Unmarshal(data, &p); err != nil {
			b.Fatal(err)
		}
	}
}
