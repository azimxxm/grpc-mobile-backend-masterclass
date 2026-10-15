// ============================================================================
// 02-server-streaming-chat: iOS Client (SwiftUI + Swift Concurrency)
// REST Polling vs gRPC Server-Streaming Solishtiruvi
// ============================================================================

import Foundation
import SwiftUI
import GRPC

// ----------------------------------------------------------------------------
// 1. AN'ANAVIY REST POLLING (Har 2 soniyada so'rov berish)
// Kamchiliklari:
// - Har 2 soniyada radiomodul (4G/5G/Wi-Fi) uyg'onadi -> Batareyani yeb tugatadi.
// - Xabar kelmagan bo'lsa ham bo'sh JSON qaytadi -> 90% trafik havoga uchadi.
// - 2 soniyagacha xabar kechikishi (Latency) bo'ladi.
// ----------------------------------------------------------------------------

@MainActor
public final class RestPollingChatViewModel: ObservableObject {
    @Published public var messages: [String] = []
    private var timer: Timer?
    private var lastMessageId: Int64 = 0
    
    public func startPolling(channelId: String) {
        timer = Timer.scheduledTimer(withTimeInterval: 2.0, repeats: true) { [weak self] _ in
            Task { [weak self] in
                await self?.pollNewMessages(channelId: channelId)
            }
        }
    }
    
    private func pollNewMessages(channelId: String) async {
        guard let url = URL(string: "https://api.example.com/channels/\(channelId)/messages?since=\(lastMessageId)") else { return }
        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            // JSON decode va massivga qo'shish...
            _ = data
        } catch {
            print("Polling error: \(error)")
        }
    }
    
    public func stop() {
        timer?.invalidate()
    }
}

// ----------------------------------------------------------------------------
// 2. gRPC SERVER-STREAMING (1 ta so'rov -> Uzluksiz jonli xabarlar oqimi)
// Yutuqlari:
// - Batareya tejaladi: Server faqat yangi xabar paydo bo'lgandagina paket uzatadi.
// - Kechikish 0 millisekund: Xabar yozilishi bilan darhol ekranga tushadi.
// - Swift-ning `AsyncSequence` / `AsyncStream` tizimi bilan tabiiy integratsiya.
// ----------------------------------------------------------------------------

@MainActor
public final class GrpcStreamingChatViewModel: ObservableObject {
    @Published public var messages: [Masterclass_V1_ChatMessage] = []
    @Published public var isConnected: Bool = false
    
    private let client: Masterclass_V1_ChatServiceAsyncClient
    private var streamTask: Task<Void, Never>?
    
    public init(channel: GRPCChannel) {
        self.client = Masterclass_V1_ChatServiceAsyncClient(channel: channel)
    }
    
    public func subscribeToChannel(channelId: String) {
        streamTask?.cancel()
        
        streamTask = Task { [weak self] in
            guard let self = self else { return }
            
            var req = Masterclass_V1_ChannelSubscriptionRequest()
            req.channelID = channelId
            req.lastMessageID = 0
            
            self.isConnected = true
            
            do {
                // Serverdan kelayotgan oqimni oddiy for-await orqali tinglaymiz:
                for try await message in self.client.subscribeChannel(req) {
                    self.messages.append(message)
                }
            } catch {
                print("gRPC Stream uzildi: \(error)")
            }
            
            self.isConnected = false
        }
    }
    
    public func disconnect() {
        streamTask?.cancel()
        streamTask = nil
    }
}
