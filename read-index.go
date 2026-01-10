package main

import (
    "encoding/binary"
    "os"
    "fmt"
)

// readIndex reads the .git/index file and returns a list of entries
func readIndex() ([]IndexEntry, error) {
    f, err := os.Open(".git/index")
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // 1. Read Header
    var header IndexHeader
    if err := binary.Read(f, binary.BigEndian, &header); err != nil {
        return nil, err
    }
    
    // Validate Signature
    if header.Signature != [4]byte{'D', 'I', 'R', 'C'} {
        return nil, fmt.Errorf("invalid index signature")
    }

    // 2. Read Entries
    entries := make([]IndexEntry, header.Entries)
    for i := 0; i < int(header.Entries); i++ {
        var entry IndexEntry
        
        // Read Fixed Fields (62 bytes)
        // We read them one by one to match the struct fields
        binary.Read(f, binary.BigEndian, &entry.CtimeSec)
        binary.Read(f, binary.BigEndian, &entry.CtimeNano)
        binary.Read(f, binary.BigEndian, &entry.MtimeSec)
        binary.Read(f, binary.BigEndian, &entry.MtimeNano)
        binary.Read(f, binary.BigEndian, &entry.Dev)
        binary.Read(f, binary.BigEndian, &entry.Ino)
        binary.Read(f, binary.BigEndian, &entry.Mode)
        binary.Read(f, binary.BigEndian, &entry.Uid)
        binary.Read(f, binary.BigEndian, &entry.Gid)
        binary.Read(f, binary.BigEndian, &entry.Size)
        binary.Read(f, binary.BigEndian, &entry.Hash)
        binary.Read(f, binary.BigEndian, &entry.Flags)

        // Read Path Name
        // The length is in the lower 12 bits of flags
        pathLen := entry.Flags & 0xFFF 
        pathBytes := make([]byte, pathLen)
        f.Read(pathBytes)
        entry.Path = string(pathBytes)

        // Read Padding
        // Logic: (62 + pathLen + padding) % 8 == 0
        // But we must consume at least 1 byte of padding.
        entryLen := 62 + pathLen
        padding := 8 - (entryLen % 8)
        
        // Skip padding bytes
        f.Seek(int64(padding), 1) // 1 means "seek relative to current"
        
        entries[i] = entry
    }
    return entries, nil
}