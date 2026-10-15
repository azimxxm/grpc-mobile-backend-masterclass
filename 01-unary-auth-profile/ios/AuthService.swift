// ============================================================================
// 01-unary-auth-profile: iOS Client (Swift 6 Concurrency)
// REST (URLSession + Codable JSON) vs gRPC (gRPC-Swift + Protobuf) Solishtiruvi
// ============================================================================

import Foundation
import GRPC
import NIO
import SwiftProtobuf

// ----------------------------------------------------------------------------
// 1. AN'ANAVIY REST USULI (URLSession + JSON)
// Kamchiliklari:
// - Matnli JSON (katta hajm, sekin parsing).
// - String kalitlar o'zgarib qolsa run-time crash xavfi bor.
// - Har safar yangi TCP ulanish yoki Keep-Alive tiqilishi.
// ----------------------------------------------------------------------------

public struct RestUserProfile: Codable, Sendable {
    public let userId: String
    public let fullName: String
    public let email: String
    public let avatarUrl: String
    public let isVerified: Bool
    public let createdAt: Int64
    
    enum CodingKeys: String, CodingKey {
        case userId = "user_id"
        case fullName = "full_name"
        case email
        case avatarUrl = "avatar_url"
        case isVerified = "is_verified"
        case createdAt = "created_at"
    }
}

public actor RestAuthManager {
    private let session = URLSession.shared
    private let baseURL = URL(string: "https://api.example.com")!
    
    public func fetchProfile(userId: String) async throws -> RestUserProfile {
        let url = baseURL.appendingPathComponent("/v1/users/\(userId)")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")
        
        let (data, response) = try await session.data(for: request)
        
        guard let httpResponse = response as? HTTPURLResponse, (200...299).contains(httpResponse.statusCode) else {
            throw URLError(.badServerResponse)
        }
        
        // Og'ir JSON decoding bosqichi:
        let decoder = JSONDecoder()
        return try decoder.decode(RestUserProfile.self, from: data)
    }
}

// ----------------------------------------------------------------------------
// 2. ZAMONAVIY gRPC USULI (gRPC-Swift + Protobuf)
// Yutuqlari:
// - Deserializatsiya 10x tezroq (Ikkilik format).
// - Type-Safety: .proto shartnomasi qonun, hech qanday xom String kalitlar yo'q!
// - HTTP/2 Multiplexing: Bitta TCP ulanish ustida parallel yuzlab so'rovlar.
// - Aniq deadline (timeout) boshqaruvi.
// ----------------------------------------------------------------------------

public actor GrpcAuthManager {
    private let client: Masterclass_V1_AuthServiceAsyncClient
    
    public init(channel: GRPCChannel) {
        self.client = Masterclass_V1_AuthServiceAsyncClient(channel: channel)
    }
    
    public func fetchProfile(userId: String) async throws -> Masterclass_V1_UserProfile {
        // 1. Shartnomaga asoslangan so'rov yaratish
        var request = Masterclass_V1_GetProfileRequest()
        request.userID = userId
        
        // 2. Qat'iy Deadline o'rnatish (masalan 3 soniya)
        var options = CallOptions()
        options.timeLimit = .timeout(.seconds(3))
        
        // 3. gRPC chaqiruvi (Type-safe, binary stream)
        let profile = try await client.getProfile(request, callOptions: options)
        return profile
    }
}
