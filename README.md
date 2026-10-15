# gRPC vs REST Masterclass: Mobile (iOS Swift) & Backend (Go)

> **Eng yuqori tezlik, minimal batareya sarfi va bank darajasidagi xavfsizlikka mo'ljallangan 5 ta real loyiha va benchmarklar to'plami.**

Ushbu repozitoriy [azimjondev.uz](https://azimjondev.uz/lessons) saytidagi **"gRPC & High-Performance System Design"** flagman kursi uchun rasmiy amaliy kodlar bazasidir.

---

## 🎯 Kurs Maqsadi: Nega Senior Dasturchilar gRPC ga O'tadi? (The AHA! Moment)

Oddiy REST API da matnli JSON, ulanishning sekinligi va run-time dagi tiqilishlar sababli yirik ilovalar (Telegram, Uber, Netflix) azob chekadi. gRPC quyidagi 5 ta fundamental ustun orqali bu muammolarni ildizi bilan yo'q qiladi:

1. **Protobuf Binary Serialization:** Matnli JSON o'rniga 10 barobar tezroq va 70% kichik ikkilik format.
2. **HTTP/2 & HTTP/3 Multiplexing:** Bitta tunnel ustida yuzlab parallel so'rovlar (TCP handshake keraksiz).
3. **Type-Safety & Single Source of Truth:** `.proto` fayl yagona qonun — Swift va Go o'rtasida API nomuvofiqligi 100% yo'qoladi.
4. **4 xil Aloqa Oqimi:** Unary, Server-Streaming, Client-Streaming, Bi-directional Streaming.
5. **Bank Darajasidagi Xavfsizlik:** mTLS (Mutual TLS), Apple Secure Enclave, SPKI Certificate Pinning va Interceptors.

---

## 📊 Katta Benchmark: JSON vs Protobuf vs FlatBuffers

| Ko'rsatkich | JSON (REST) | Protobuf (gRPC) | FlatBuffers (Zero-Copy) |
| :--- | :--- | :--- | :--- |
| **Deserializatsiya (Decoding)** | ~120 - 250 ms *(Sekin)* | ~15 - 25 ms *(Tez)* | **~0.02 ms** ⚡ *(Bir zumda)* |
| **Qo'shimcha RAM Allokatsiyasi** | Juda katta (Obyektlar daraxti) | O'rtacha | **0 bayt** (In-place pointer) |
| **Paket Hajmi (Wire Size)** | 1.2 MB (100%) | 250 KB (20%) | 320 KB (25%) |
| **CPU va Batareya Sarfi** | Yuqori (Telefon qiziydi) | Minimal | Eng past |

---

## 🏗 5 Ta Real-World Demo Loyihalar

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
- **Mavzu:** Bank o'tkazmalari va to'lov tizimlari xavfsizligi.
- **Taqqoslash:** Qo'lda yozilgan URLSession headerlari vs `mTLS + Client/Server Interceptors`.
- **Yutuq:** Idempotency tekshiruvi, Dynamic Nonce, replay-attack dan to'liq himoyalangan avtomatlashtirilgan middleware.

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
