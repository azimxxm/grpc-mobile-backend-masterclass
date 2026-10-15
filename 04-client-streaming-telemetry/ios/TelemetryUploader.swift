// ============================================================================
// 04-client-streaming-telemetry: iOS Client
// Client-Streaming RPC: Phone streams telemetry events -> Server returns summary
// ============================================================================

import Foundation
import GRPC

public actor TelemetryUploader {
    private let client: Masterclass_V1_TelemetryServiceAsyncClient
    
    public init(channel: GRPCChannel) {
        self.client = Masterclass_V1_TelemetryServiceAsyncClient(channel: channel)
    }
    
    public func uploadBufferedEvents(_ events: [Masterclass_V1_MetricEvent]) async throws -> Masterclass_V1_UploadSummary {
        // 1. Client-Streaming chaqiruvini ochamiz
        let stream = client.uploadBatchMetrics()
        
        // 2. Har bir voqeani oqim orqali uzatamiz (katta xotira allokatsiyalarisiz):
        for event in events {
            try await stream.requestStream.send(event)
        }
        
        // 3. Oqimni tugatamiz
        try await stream.requestStream.finish()
        
        // 4. Serverdan yakuniy xulosa javobini qabul qilamiz:
        let summary = try await stream.response
        return summary
    }
}
