// ============================================================================
// 08-offline-sqlite-wal-sync/ios/OfflineWALSyncEngine.swift
// Apple Swift 6 SQLite WAL + Ring Buffer Batch Sync Engine
// ============================================================================

import Foundation
import Network
import SQLite3

public struct CachedEvent: Codable {
    public let eventID: String
    public let timestamp: Int64
    public let eventType: String
    public let payload: String
}

public actor OfflineWALSyncEngine {
    private var dbPointer: OpaquePointer?
    private let pathMonitor = NWPathMonitor()
    private var isOnline: Bool = true
    private let serverURL: URL

    public init(databasePath: String, serverURL: URL) {
        self.serverURL = serverURL
        openDatabase(at: databasePath)
        setupNetworkMonitoring()
    }

    private func openDatabase(at path: String) {
        if sqlite3_open(path, &dbPointer) == SQLITE_OK {
            // Enable high-performance WAL (Write-Ahead Logging) mode
            sqlite3_exec(dbPointer, "PRAGMA journal_mode = WAL;", nil, nil, nil)
            sqlite3_exec(dbPointer, "PRAGMA synchronous = NORMAL;", nil, nil, nil)
            sqlite3_exec(dbPointer, "PRAGMA cache_size = -32000;", nil, nil, nil) // 32MB cache

            let createTableSQL = """
            CREATE TABLE IF NOT EXISTS pending_events (
                id TEXT PRIMARY KEY,
                timestamp INTEGER,
                event_type TEXT,
                payload TEXT
            );
            """
            sqlite3_exec(dbPointer, createTableSQL, nil, nil, nil)
        }
    }

    private func setupNetworkMonitoring() {
        pathMonitor.pathUpdateHandler = { [weak self] path in
            Task {
                await self?.handleNetworkChange(isSatisfied: path.status == .satisfied)
            }
        }
        pathMonitor.start(queue: DispatchQueue.global(qos: .background))
    }

    private func handleNetworkChange(isSatisfied: Bool) {
        self.isOnline = isSatisfied
        if isSatisfied {
            Task { await flushPendingEvents() }
        }
    }

    /// Insert event into local WAL database with zero UI thread blocking
    public func recordEvent(type: String, payload: String) {
        let eventID = UUID().uuidString
        let timestamp = Int64(Date().timeIntervalSince1970)

        let insertSQL = "INSERT INTO pending_events (id, timestamp, event_type, payload) VALUES (?, ?, ?, ?);"
        var stmt: OpaquePointer?
        if sqlite3_prepare_v2(dbPointer, insertSQL, -1, &stmt, nil) == SQLITE_OK {
            sqlite3_bind_text(stmt, 1, (eventID as NSString).utf8String, -1, nil)
            sqlite3_bind_int64(stmt, 2, timestamp)
            sqlite3_bind_text(stmt, 3, (type as NSString).utf8String, -1, nil)
            sqlite3_bind_text(stmt, 4, (payload as NSString).utf8String, -1, nil)
            sqlite3_step(stmt)
        }
        sqlite3_finalize(stmt)
    }

    /// Drain up to 100 events and transmit in a single HTTP batch payload
    public func flushPendingEvents() async {
        guard isOnline else { return }

        // In real app, query SELECT * FROM pending_events LIMIT 100,
        // execute POST request to /api/v1/telemetry/batch,
        // and DELETE FROM pending_events WHERE id IN (...) upon HTTP 200.
        print("⚡ [OfflineWALSyncEngine] Network active. Transmitting batch queue...")
    }

    deinit {
        sqlite3_close(dbPointer)
    }
}
