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


    var header IndexHeader
    if err := binary.Read(f, binary.BigEndian, &header); err != nil {
        return nil, err
    }
    

    if header.Signature != [4]byte{'D', 'I', 'R', 'C'} {
        return nil, fmt.Errorf("invalid index signature")
    }


    entries := make([]IndexEntry, header.Entries)
    for i := 0; i < int(header.Entries); i++ {
        var entry IndexEntry
        

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


        pathLen := entry.Flags & 0xFFF 
        pathBytes := make([]byte, pathLen)
        f.Read(pathBytes)
        entry.Path = string(pathBytes)


        entryLen := 62 + pathLen
        padding := 8 - (entryLen % 8)
        

        f.Seek(int64(padding), 1) 
        
        entries[i] = entry
    }
    return entries, nil
}