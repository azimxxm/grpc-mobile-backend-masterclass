// ============================================================================
// 06-flatbuffers-zero-copy/backend/main.go
// High-Speed FlatBuffers Binary Server in Go
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"

	flatbuffers "github.com/google/flatbuffers/go"
)

// Standalone builder demonstration without external code-gen dependencies
func buildEntityBinary(id uint64, name string, hp int16, mana int16, posX, posY, posZ float32) []byte {
	b := flatbuffers.NewBuilder(1024)

	// String serialized
	nameOffset := b.CreateString(name)

	// Start object
	b.StartObject(7)
	b.PrependUint64Slot(0, id, 0)
	b.PrependUOffsetTRelativeSlot(1, nameOffset, 0)
	b.PrependInt16Slot(3, hp, 100)
	b.PrependInt16Slot(4, mana, 50)
	entityOffset := b.EndObject()

	b.Finish(entityOffset)
	return b.FinishedBytes()
}

func main() {
	http.HandleFunc("/api/v1/entity/binary", func(w http.ResponseWriter, r *http.Request) {
		binaryData := buildEntityBinary(99991, "Shadow Dragon", 2500, 800, 120.5, 45.0, -89.2)

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("X-Serialization-Engine", "Google-FlatBuffers-Zero-Copy")
		w.Write(binaryData)
	})

	fmt.Println("🚀 FlatBuffers Zero-Copy Go Server listening on :8086...")
	log.Fatal(http.ListenAndServe(":8086", nil))
}
