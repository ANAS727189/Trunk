package main

// The Header
type IndexHeader struct {
    Signature [4]byte // "DIRC"
    Version   uint32  // 2
    Entries   uint32  // Count of files
}

// The Entry
type IndexEntry struct {
    CtimeSec  uint32
    CtimeNano uint32
    MtimeSec  uint32
    MtimeNano uint32
    Dev       uint32
    Ino       uint32
    Mode      uint32
    Uid       uint32
    Gid       uint32
    Size      uint32
    Hash      [20]byte // The Blob Hash (Raw bytes, not hex string)
    Flags     uint16
    Path      string
}