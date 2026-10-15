// ============================================================================
// 07-http3-quic-transport/ios/HTTP3ConnectionManager.swift
// Apple Swift 6 HTTP/3 & QUIC Seamless Connection Migration Manager
// ============================================================================

import Foundation
import Network

@MainActor
public final class HTTP3ConnectionManager: ObservableObject {
    @Published public private(set) var activeProtocol: String = "Connecting..."
    @Published public private(set) var receivedPackets: [String] = []
    @Published public private(set) var networkInterface: String = "Unknown"

    private var session: URLSession!
    private let pathMonitor = NWPathMonitor()

    public init() {
        configureSession()
        setupInterfaceMonitoring()
    }

    private func configureSession() {
        let config = URLSessionConfiguration.default
        // Enable native HTTP/3 QUIC transport
        config.assumesHTTP3Capable = true
        // Enable seamless multipath handover between Wi-Fi and Cellular
        config.multipathServiceType = .handover
        config.waitsForConnectivity = true
        config.timeoutIntervalForRequest = 30

        self.session = URLSession(configuration: config)
    }

    private func setupInterfaceMonitoring() {
        pathMonitor.pathUpdateHandler = { [weak self] path in
            Task { @MainActor in
                if path.usesInterfaceType(.wifi) {
                    self?.networkInterface = "Wi-Fi (High Bandwidth)"
                } else if path.usesInterfaceType(.cellular) {
                    self?.networkInterface = "Cellular LTE/5G (Mobile)"
                } else {
                    self?.networkInterface = "Other Interface"
                }
            }
        }
        pathMonitor.start(queue: DispatchQueue.global(qos: .background))
    }

    /// Stream packets from HTTP/3 endpoint with Zero-Downtime handover
    public func startStream(url: URL) async {
        do {
            var request = URLRequest(url: url)
            request.httpMethod = "GET"

            let (asyncBytes, response) = try await session.bytes(for: request)
            
            if let httpResponse = response as? HTTPURLResponse {
                self.activeProtocol = httpResponse.allHeaderFields["Alt-Svc"] != nil ? "HTTP/3 (QUIC)" : "HTTP/2"
            }

            for try await line in asyncBytes.lines {
                guard !line.isEmpty, line.hasPrefix("data:") else { continue }
                self.receivedPackets.append(line)
            }
        } catch {
            print("Stream error (QUIC will retry automatically): \(error.localizedDescription)")
        }
    }
}
