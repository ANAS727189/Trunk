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

    entries, _ := readIndex() 
    if entries == nil {
        entries = []IndexEntry{}
    }


    info, err := os.Stat(filename)
    if err != nil {
        log.Fatal("File not found:", err)
    }
    stat := info.Sys().(*syscall.Stat_t)


    hashBytes := hashObject(filename)


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


    idx := -1
    for i, e := range entries {
        if e.Path == filename {
            idx = i
            break
        }
    }

    if idx != -1 {
        entries[idx] = newEntry
    } else {
        entries = append(entries, newEntry)
    }


    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Path < entries[j].Path
    })


    f, err := os.Create(".git/index")
    if err != nil {
        log.Fatal("Could not create index:", err)
    }
    defer f.Close()


    header := IndexHeader{
        Signature: [4]byte{'D', 'I', 'R', 'C'},
        Version:   2,
        Entries:   uint32(len(entries)), 
    }
    binary.Write(f, binary.BigEndian, header)


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