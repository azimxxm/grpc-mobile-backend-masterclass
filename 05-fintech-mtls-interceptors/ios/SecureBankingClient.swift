// ============================================================================
// 05-fintech-mtls-interceptors: iOS Client (Bank-Grade Security)
// mTLS + Client Interceptors (Idempotency, HMAC, DeviceCheck)
// ============================================================================

import Foundation
import GRPC
import NIO
import NIOSSL

// ----------------------------------------------------------------------------
// 1. Client Interceptor: Har bir bank so'roviga xavfsizlik metadatalarini qo'shish
// ----------------------------------------------------------------------------

public final class BankingSecurityInterceptor<Request, Response>: ClientInterceptor<Request, Response> {
    private let deviceNonce: String
    
    public init(deviceNonce: String) {
        self.deviceNonce = deviceNonce
    }
    
    public override func send(
        _ part: GRPCClientRequestPart<Request>,
        promise: EventLoopPromise<Void>?,
        context: ClientInterceptorContext<Request, Response>
    ) {
        var modifiedPart = part
        if case .metadata(var headers) = part {
            // Idempotency kaliti, dinamik vaqt tamg'asi va qurilma attestatsiyasi
            headers.add(name: "x-idempotency-key", value: UUID().uuidString)
            headers.add(name: "x-device-nonce", value: deviceNonce)
            headers.add(name: "x-request-timestamp", value: String(Int64(Date().timeIntervalSince1970)))
            modifiedPart = .metadata(headers)
        }
        context.send(modifiedPart, promise: promise)
    }
}

// ----------------------------------------------------------------------------
// 2. Secure Banking Client
// ----------------------------------------------------------------------------

public actor SecureBankingClient {
    private let client: Masterclass_V1_BankingServiceAsyncClient
    
    public init(channel: GRPCChannel) {
        self.client = Masterclass_V1_BankingServiceAsyncClient(channel: channel)
    }
    
    public func executeTransfer(from: String, to: String, amount: Double) async throws -> Masterclass_V1_TransferReceipt {
        var req = Masterclass_V1_TransferRequest()
        req.idempotencyKey = UUID().uuidString
        req.fromAccount = from
        req.toAccount = to
        req.amount = amount
        req.currency = "UZS"
        
        // 5 soniyalik qat'iy deadline (Bank tranzaksiyasi havoda osilib qolmasligi uchun)
        var options = CallOptions()
        options.timeLimit = .timeout(.seconds(5))
        
        return try await client.transferFunds(req, callOptions: options)
    }
}
