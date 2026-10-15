// ============================================================================
// 03-bidirectional-geotracking: iOS Client (Uber-style Live Radar)
// Bi-directional Streaming RPC: Client sends GPS -> Server streams nearby drivers
// ============================================================================

import Foundation
import CoreLocation
import GRPC
import SwiftUI

@MainActor
public final class LiveRadarViewModel: ObservableObject {
    @Published public var nearbyDrivers: [Masterclass_V1_DriverInfo] = []
    @Published public var isTrackingActive: Bool = false
    
    private let client: Masterclass_V1_TrackingServiceAsyncClient
    private var trackingTask: Task<Void, Never>?
    
    public init(channel: GRPCChannel) {
        self.client = Masterclass_V1_TrackingServiceAsyncClient(channel: channel)
    }
    
    public func startLiveRadar(clientId: String) {
        trackingTask?.cancel()
        
        trackingTask = Task { [weak self] in
            guard let self = self else { return }
            self.isTrackingActive = true
            
            // 1. Bi-directional streamni ochamiz
            let stream = self.client.liveLocationStream()
            
            // 2. Parallel ravishda serverdan keluvchi haydovchilar radarini o'qish:
            let responseTask = Task {
                do {
                    for try await radar in stream.responseStream {
                        self.nearbyDrivers = radar.drivers
                    }
                } catch {
                    print("Radar oqimi to'xtadi: \(error)")
                }
            }
            
            // 3. Telefon harakatlanganda o'z GPS koordinatalarimizni serverga yuborish:
            var currentLat = 41.2995
            var currentLon = 69.2401
            
            for _ in 1...60 { // Simulyatsiya: har 1 soniyada koordinata uzatish
                if Task.isCancelled { break }
                
                var update = Masterclass_V1_ClientLocationUpdate()
                update.clientID = clientId
                update.latitude = currentLat
                update.longitude = currentLon
                update.timestamp = Int64(Date().timeIntervalSince1970)
                
                do {
                    try await stream.requestStream.send(update)
                } catch {
                    print("GPS yuborishda xato: \(error)")
                    break
                }
                
                currentLat += 0.0001
                currentLon += 0.0001
                try? await Task.sleep(nanoseconds: 1_000_000_000)
            }
            
            // Oqimni yopish
            try? await stream.requestStream.finish()
            _ = await responseTask.result
            self.isTrackingActive = false
        }
    }
    
    public func stop() {
        trackingTask?.cancel()
        trackingTask = nil
    }
}
