package tttandme

import (
    "github.com/edsrzf/mmap-go"
)

func consume(delim byte, i uint64, m mmap.MMap) uint64 {
    for m[i] != delim {
        i++
    }
    return i
}
