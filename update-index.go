package main

import (
    "encoding/binary"
    "fmt"
    "log"
    "os"
    "sort" 
    "syscall"
)

func updateIndex(filename string) {
    // 1. Read Existing Index
    entries, _ := readIndex() 
    if entries == nil {
        entries = []IndexEntry{}
    }

    // 2. Stat the new file
    info, err := os.Stat(filename)
    if err != nil {
        log.Fatal("File not found:", err)
    }
    stat := info.Sys().(*syscall.Stat_t)

    // 3. Create Blob and Get Hash
    hashBytes := hashObject(filename)

    // 4. Create the New Entry
    newEntry := IndexEntry{
        CtimeSec:  uint32(stat.Ctim.Sec),
        CtimeNano: uint32(stat.Ctim.Nsec),
        MtimeSec:  uint32(stat.Mtim.Sec),
        MtimeNano: uint32(stat.Mtim.Nsec),
        Dev:       uint32(stat.Dev),
        Ino:       uint32(stat.Ino),
        Mode:      0x81A4, // 100644
        Uid:       uint32(stat.Uid),
        Gid:       uint32(stat.Gid),
        Size:      uint32(info.Size()),
        Flags:     uint16(len(filename)),
        Path:      filename,
    }
    copy(newEntry.Hash[:], hashBytes)

    // 5. Update or Append
    // We check if the file is already in our list
    idx := -1
    for i, e := range entries {
        if e.Path == filename {
            idx = i
            break
        }
    }

    if idx != -1 {
        // Update existing entry
        entries[idx] = newEntry
    } else {
        // Append new entry
        entries = append(entries, newEntry)
    }

    // 6. SORT THE ENTRIES (Crucial for Git!)
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Path < entries[j].Path
    })

    // 7. Write the Index File
    f, err := os.Create(".git/index")
    if err != nil {
        log.Fatal("Could not create index:", err)
    }
    defer f.Close()

    // Write Header
    header := IndexHeader{
        Signature: [4]byte{'D', 'I', 'R', 'C'},
        Version:   2,
        Entries:   uint32(len(entries)), // Now reflects total count
    }
    binary.Write(f, binary.BigEndian, header)

    // Write Entries
    for _, entry := range entries {
        binary.Write(f, binary.BigEndian, entry.CtimeSec)
        binary.Write(f, binary.BigEndian, entry.CtimeNano)
        binary.Write(f, binary.BigEndian, entry.MtimeSec)
        binary.Write(f, binary.BigEndian, entry.MtimeNano)
        binary.Write(f, binary.BigEndian, entry.Dev)
        binary.Write(f, binary.BigEndian, entry.Ino)
        binary.Write(f, binary.BigEndian, entry.Mode)
        binary.Write(f, binary.BigEndian, entry.Uid)
        binary.Write(f, binary.BigEndian, entry.Gid)
        binary.Write(f, binary.BigEndian, entry.Size)
        binary.Write(f, binary.BigEndian, entry.Hash)
        binary.Write(f, binary.BigEndian, entry.Flags)

        f.Write([]byte(entry.Path))

        entryLen := 62 + len(entry.Path)
        padding := 8 - (entryLen % 8)
        f.Write(make([]byte, padding))
    }

    fmt.Printf("Index updated. Tracking %d files.\n", len(entries))
}