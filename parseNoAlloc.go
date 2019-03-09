package tttandme

import (
    "os"
    "fmt"
    "bytes"
    "strconv"
    "reflect"
    "unsafe"

    "github.com/edsrzf/mmap-go"
)

const (
    numMaps    = 24
    partsBound = 0x10000000000000000 / numMaps
)

type genomeNoMemory struct {
    m         mmap.MMap
    y         uint8
    idSmallI  []map[uint64]uint64
    idSmallRS []map[uint64]uint64
    idLargeI  map[string]uint64
    idLargeRS map[string]uint64
}

// This parser won't keep the parsed contents in memory.
// Using a parsed Genome after the parser has closed is
// undefined behavior, which will likely result in a runtime error.
func NewParserNoMem() GenomeParser {
    return &genomeNoMemory{}
}

func (g *genomeNoMemory) Open(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    m, err := mmap.Map(f, mmap.RDONLY, 0)
    if err != nil {
        f.Close()
        return err
    }
    g.m = m

    return nil
}

func (g *genomeNoMemory) Parse() (Genome, error) {
    g.y = 0
    g.initMaps()
    lim := uint64(len(g.m))

    for i := uint64(0); i < lim; {
        iNewline := consume('\n', i, g.m)

        // comment, do nothing
        if g.m[i] == '#' {
            i = iNewline + 1
            continue
        }

        // check for male
        if g.y == 0 {
            s := []byte(g.m[i:iNewline])
            if bytes.IndexByte(s, 'Y') > 0 && bytes.IndexByte(s, '-') < 0 {
                g.y = 2
            }
        }

        // find tab and save in rsid map
        iTab := consume('\t', i, g.m)
        id := toString(g.m[i:iTab])
        g.setRSID(id, iTab+1)

        i = iNewline + 1
    }

    g.y++

    return g, nil
}

func (g *genomeNoMemory) RSID(id string) *SNP {
    i, ok := g.getRSID(id)
    if !ok {
        return nil
    }

    // find all token positions
    iTabPos := consume('\t', i, g.m)
    iTabGen := consume('\t', iTabPos+1, g.m)
    iNewline := consume('\n', iTabGen+1, g.m)

    // retrieve tokens
    _chrxm := toString(g.m[i:iTabPos])
    _pos := toString(g.m[iTabPos+1:iTabGen])
    _genotp := toString(g.m[iTabGen+1:iNewline])

    // convert chromossome
    var chrxm Chromosome

    switch {
    default:
        return nil
    case _chrxm[0] >= '0' && _chrxm[0] <= '9':
        i, err := strconv.Atoi(_chrxm)
        if err != nil {
            return nil
        }
        chrxm = Chromosome(i)
    case _chrxm == "X":
        chrxm = CHR_X
    case _chrxm == "Y":
        chrxm = CHR_Y
    case _chrxm == "MT":
        chrxm = CHR_MT
    }

    // convert position
    var pos Position

    i, err := strconv.ParseUint(_pos, 10, 64)
    if err != nil {
        return nil
    }
    pos = Position(i)

    return &SNP{
        Chromosome: chrxm,
        Position: pos,
        Genotype: Geno(_genotp),
    }
}

func (g *genomeNoMemory) Iter(f func(string) bool) bool {
    for id := range g.idLargeI {
        if !f(id) {
            return false
        }
    }
    for id := range g.idLargeRS {
        if !f(id) {
            return false
        }
    }
    for i := 0; i < numMaps; i++ {
        for id := range g.idSmallI[i] {
            if !f(fmt.Sprintf("i%d", id)) {
                return false
            }
        }
    }
    for i := 0; i < numMaps; i++ {
        for id := range g.idSmallRS[i] {
            if !f(fmt.Sprintf("i%d", id)) {
                return false
            }
        }
    }
    return true
}

func (g *genomeNoMemory) HasY() bool {
    return g.y > 1
}

func (g *genomeNoMemory) Close() error {
    g.freeMaps()
    return g.m.Unmap()
}

func (g *genomeNoMemory) initMaps() {
    g.idSmallI = make([]map[uint64]uint64, numMaps)
    g.idSmallRS = make([]map[uint64]uint64, numMaps)
    for i := 0; i < numMaps; i++ {
        g.idSmallI[i] = make(map[uint64]uint64)
        g.idSmallRS[i] = make(map[uint64]uint64)
    }
    g.idLargeI = make(map[string]uint64)
    g.idLargeRS = make(map[string]uint64)
}

func (g *genomeNoMemory) freeMaps() {
    g.idSmallI = nil
    g.idSmallRS = nil
    g.idLargeI = nil
    g.idLargeRS = nil
}

func (g *genomeNoMemory) getRSID(rsid string) (i uint64, ok bool) {
    switch {
    case rsid[0] == 'r' || rsid[0] == 'R':
        rsid = rsid[2:]
        if len(rsid) > 8 {
            i, ok = g.idLargeRS[rsid]
            return
        }
        key := getRSIDKey(&rsid)
        i, ok = g.idSmallRS[mapId(key)][key]
    case rsid[0] == 'i' || rsid[0] == 'I':
        rsid = rsid[1:]
        if len(rsid) > 8 {
            i, ok = g.idLargeI[rsid]
            return
        }
        key := getRSIDKey(&rsid)
        i, ok = g.idSmallI[mapId(key)][key]
    }
    return
}

func (g *genomeNoMemory) setRSID(rsid string, i uint64) {
    switch {
    case rsid[0] == 'r' || rsid[0] == 'R':
        rsid = rsid[2:]
        if len(rsid) > 8 {
            g.idLargeRS[rsid] = i
            return
        }
        key := getRSIDKey(&rsid)
        g.idSmallRS[mapId(key)][key] = i
    case rsid[0] == 'i' || rsid[0] == 'I':
        rsid = rsid[1:]
        if len(rsid) > 8 {
            g.idLargeI[rsid] = i
            return
        }
        key := getRSIDKey(&rsid)
        g.idSmallI[mapId(key)][key] = i
    }
}

func getRSIDKey(rsid *string) uint64 {
    s := (*reflect.StringHeader)(unsafe.Pointer(rsid))
    key := (*uint64)(unsafe.Pointer(s.Data))

    // remove memory garbage
    switch s.Len {
    default:
        return 0
    case 1:
        return *key & 0xff
    case 2:
        return *key & 0xffff
    case 3:
        return *key & 0xffffff
    case 4:
        return *key & 0xffffffff
    case 5:
        return *key & 0xffffffffff
    case 6:
        return *key & 0xffffffffffff
    case 7:
        return *key & 0xffffffffffffff
    case 8:
        return *key
    }
}

func mapId(key uint64) uint64 {
    return key / partsBound
}

func toString(s mmap.MMap) string {
    sh := reflect.StringHeader{
        Data: uintptr(unsafe.Pointer(&s[0])),
        Len: len(s),
    }
    return *(*string)(unsafe.Pointer(&sh))
}
