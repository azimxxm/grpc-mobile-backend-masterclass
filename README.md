# gRPC & High-Performance System Design Masterclass: iOS (Swift 6) & Go

> **Eng yuqori tezlik, minimal batareya sarfi, Zero-Copy xotira va bank darajasidagi xavfsizlikka mo'ljallangan 8 ta real loyiha va benchmarklar to'plami.**

Ushbu repozitoriy [azimjondev.uz](https://azimjondev.uz/lessons) saytidagi **"gRPC & High-Performance System Design"** flagman masterclass kursi uchun rasmiy amaliy kodlar bazasidir.

---

## 📚 8 Ta Flagman Darslik (azimjondev.uz)

1. **`01_the_cost_of_json.md`**: Nega JSON yirik loyihalarda o'ladi? String parsing, Heap allocations va benchmarklar.
2. **`02_protobuf_and_grpc.md`**: Protobuf v3 + gRPC to'liq qo'llanma (Swift + Go). Varint, Field tags va HTTP/2 multiplexing.
3. **`03_flatbuffers_zero_copy.md`**: FlatBuffers: 0 ms parslash va xotira mo'jizasi (Zero-Copy VTable pointers).
4. **`04_telegram_mtproto_tl.md`**: Telegram TL va MTProto qanday ishlaydi? Nikolay Durov algoritmlari va In-Memory engine.
5. **`05_http3_quic_transport.md`**: UDP asosidagi HTTP/3: Head-of-Line blocking barhami va Connection Migration.
6. **`06_mobile_security_mtls_pinning.md`**: mTLS, Certificate Pinning, Apple Secure Enclave va Dynamic HMAC.
7. **`07_offline_first_caching.md`**: SQLite WAL, Ring Buffer va Batch Sync Engine.
8. **`08_real_world_projects_showcase.md`**: 8 ta to'liq demo loyiha kodlari va masterclass xulosasi.

---

## 🏗 8 Ta Real-World Demo Loyihalar (Kodlar Bazasi)

Har bir papkada to'liq ishlovchi **iOS Swift 6 Client** va **Go Backend** kodlari mavjud:

### 📁 [`01-unary-auth-profile/`](01-unary-auth-profile/)
- **Mavzu:** Klassik Auth va Foydalanuvchi Profili yuklash.
- **Taqqoslash:** `URLSession + Codable JSON` vs `gRPC-Swift + Protobuf Unary RPC`.
- **Yutuq:** API xavfsizligi, 70% kamroq trafik, qat'iy Deadline boshqaruvi.

### 📁 [`02-server-streaming-chat/`](02-server-streaming-chat/)
- **Mavzu:** Jonli xabarlar va yangiliklar tasmasi (Telegram / Slack).
- **Taqqoslash:** REST Polling (har 2 soniyada so'rov) vs `gRPC Server-Streaming`.
- **Yutuq:** 0 ms kechikish, radiomodul tinch holatda turgani sababli batareya 3x kamroq sarflanadi, Swift `AsyncStream` integratsiyasi.

### 📁 [`03-bidirectional-geotracking/`](03-bidirectional-geotracking/)
- **Mavzu:** Jonli mashinalar radari (Uber / Yandex Go).
- **Taqqoslash:** Har sekundda `POST /location` yuborish vs `gRPC Bi-directional Stream`.
- **Yutuq:** Mijoz o'z koordinatasini uzatib turadi, server esa ayni vaqtda uning atrofidagi haydovchilarni qaytaradi.

### 📁 [`04-client-streaming-telemetry/`](04-client-streaming-telemetry/)
- **Mavzu:** Katta hajmdagi reklama va xatoliklar telemetriyasi (NovagramAds / Datadog).
- **Taqqoslash:** 100 ta voqeani ulkan JSON array qilib yuborish vs `gRPC Client-Streaming`.
- **Yutuq:** Xotirani to'ldirmasdan voqealar birma-bir oqiziladi, server yakuniy hisobotni qaytaradi.

### 📁 [`05-fintech-mtls-interceptors/`](05-fintech-mtls-interceptors/)
- **Mavzu:** Bank o'tkazmalari va to'lov tizimlari xavfsizligi (Payme / Click / Bank).
- **Taqqoslash:** Qo'lda yozilgan URLSession headerlari vs `mTLS + Client/Server Interceptors`.
- **Yutuq:** Idempotency tekshiruvi, Dynamic Nonce, replay-attack dan to'liq himoyalangan avtomatlashtirilgan middleware.

### 📁 [`06-flatbuffers-zero-copy/`](06-flatbuffers-zero-copy/)
- **Mavzu:** O'yinlar va 3D grafika uchun Google FlatBuffers Zero-Copy Reader.
- **Taqqoslash:** Protobuf Deserialization (18 ms) vs FlatBuffers Direct ByteBuffer Pointer (**0.02 ms**).
- **Yutuq:** 0 bayt yangi RAM allokatsiyasi, instant O(1) field lookup.

### 📁 [`07-http3-quic-transport/`](07-http3-quic-transport/)
- **Mavzu:** UDP ustidagi HTTP/3 QUIC va Mobil Connection Migration.
- **Taqqoslash:** TCP Connection Reset on Wi-Fi/LTE drop vs QUIC Seamless CID Migration.
- **Yutuq:** `Network.framework` handover orqali 0-sekundlik uzilishsiz ma'lumot oqimi.

### 📁 [`08-offline-sqlite-wal-sync/`](08-offline-sqlite-wal-sync/)
- **Mavzu:** Offline-First SQLite WAL rejimi va High-Throughput Batch Ingestion.
- **Taqqoslash:** Har bir event uchun alohida HTTP so'rov vs Local WAL Buffer + Single Transaction Ingest.
- **Yutuq:** Readers never block Writers (`PRAGMA journal_mode=WAL`), batareya tejamkorligi.

---

## 📊 Katta Benchmark: JSON vs Protobuf vs FlatBuffers

| Ko'rsatkich | JSON (REST) | Protobuf (gRPC) | FlatBuffers (Zero-Copy) |
| :--- | :--- | :--- | :--- |
| **Deserializatsiya (Decoding)** | ~120 - 250 ms *(Sekin)* | ~15 - 25 ms *(Tez)* | **~0.02 ms** ⚡ *(Bir zumda)* |
| **Qo'shimcha RAM Allokatsiyasi** | Juda katta (Obyektlar daraxti) | O'rtacha | **0 bayt** (In-place pointer) |
| **Paket Hajmi (Wire Size)** | 1.2 MB (100%) | 250 KB (20%) | 320 KB (25%) |
| **CPU va Batareya Sarfi** | Yuqori (Telefon qiziydi) | Minimal | Eng past |

---

## 🛠 Qanday Ishga Tushiriladi?

### 1. Protobuf Kompilyatsiyasi (Agar `.proto` ga o'zgartirish kiritilsa):
```bash
# Go uchun generatsiya:
protoc --go_out=gen/go --go-grpc_out=gen/go proto/services.proto

# Swift uchun generatsiya:
protoc --swift_out=gen/swift --grpc-swift_out=gen/swift proto/services.proto
```

### 2. Go Backend-ni Ishga Tushirish:
```bash
# Har bir loyiha papkasiga kirib:
cd 01-unary-auth-profile/backend
go run main.go
```

### 3. Benchmarklarni Sinab Ko'rish:
```bash
cd benchmarks
go test -bench=. -benchmem
```

---

## 👨‍💻 Muallif & Darsliklar
- Muallif: **Azimjon Po'latov** ([azimjondev.uz](https://azimjondev.uz))
- Savol va takliflar uchun: [Telegram](https://t.me/azimxxm)
