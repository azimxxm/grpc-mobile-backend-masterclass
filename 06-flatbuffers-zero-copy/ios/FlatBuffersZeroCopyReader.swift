// ============================================================================
// 06-flatbuffers-zero-copy/ios/FlatBuffersZeroCopyReader.swift
// Apple Swift 6 Zero-Copy Binary Memory Reader
// ============================================================================

import Foundation

/// Zero-Copy Binary Reader that parses FlatBuffers memory offsets directly
/// without heap allocations or object instantiation overhead.
public final class FlatBuffersZeroCopyReader {
    private let rawData: Data

    public init(binaryData: Data) {
        self.rawData = binaryData
    }

    /// Reads an integer directly from the byte array using raw pointer offsets (O(1))
    public func readUInt64(at offset: Int) -> UInt64? {
        guard offset + MemoryLayout<UInt64>.size <= rawData.count else { return nil }
        return rawData.withUnsafeBytes { rawBufferPointer in
            rawBufferPointer.loadUnaligned(fromByteOffset: offset, as: UInt64.self)
        }
    }

    public func readInt16(at offset: Int) -> Int16? {
        guard offset + MemoryLayout<Int16>.size <= rawData.count else { return nil }
        return rawData.withUnsafeBytes { rawBufferPointer in
            rawBufferPointer.loadUnaligned(fromByteOffset: offset, as: Int16.self)
        }
    }

    public func readFloat32(at offset: Int) -> Float32? {
        guard offset + MemoryLayout<Float32>.size <= rawData.count else { return nil }
        return rawData.withUnsafeBytes { rawBufferPointer in
            rawBufferPointer.loadUnaligned(fromByteOffset: offset, as: Float32.self)
        }
    }

    /// Demonstration of instant O(1) field inspection
    public func benchmarkDirectRead() -> (id: UInt64, hp: Int16, timeNanoseconds: UInt64) {
        let start = DispatchTime.now().uptimeNanoseconds
        
        // Root offset pointer in FlatBuffers is at offset 0 (4 bytes)
        let rootOffset = rawData.withUnsafeBytes { $0.loadUnaligned(fromByteOffset: 0, as: UInt32.self) }
        
        let id = readUInt64(at: Int(rootOffset)) ?? 0
        let hp = readInt16(at: Int(rootOffset) + 8) ?? 0

        let end = DispatchTime.now().uptimeNanoseconds
        return (id, hp, end - start)
    }
}
