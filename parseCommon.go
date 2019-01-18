package tttandme

import (
    "github.com/edsrzf/mmap-go"
)

func consume(delim byte, i uint64, m mmap.MMap) uint64 {
    for {
        if m[i] == delim {
            return i
        }
        i++
        if m[i] == delim {
            return i
        }
        i++
        if m[i] == delim {
            return i
        }
        i++
        if m[i] == delim {
            return i
        }
        i++
    }
}
