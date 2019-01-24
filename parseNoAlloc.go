package tttandme

import (
    "os"
    "bytes"
    "strconv"

    "github.com/edsrzf/mmap-go"
)

type genomeNoMemory struct {
    m    mmap.MMap
    rsid map[string]uint64
    y    uint8
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

    rsid := make(map[string]uint64)
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
        id := string(g.m[i:iTab])
        rsid[id] = iTab+1

        i = iNewline + 1
    }

    g.y++
    g.rsid = rsid

    return g, nil
}

func (g *genomeNoMemory) RSID(id string) *SNP {
    i, ok := g.rsid[id]
    if !ok {
        return nil
    }

    // find all token positions
    iTabPos := consume('\t', i, g.m)
    iTabGen := consume('\t', iTabPos+1, g.m)
    iNewline := consume('\n', iTabGen+1, g.m)

    // retrieve tokens
    _chrxm := string(g.m[i:iTabPos])
    _pos := string(g.m[iTabPos+1:iTabGen])
    _genotp := string(g.m[iTabGen+1:iNewline])

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
        Genotype: GenotypeStr(_genotp),
    }
}

func (g *genomeNoMemory) Iter(f func(string) bool) bool {
    for id := range g.rsid {
        if !f(id) {
            return false
        }
    }
    return true
}

func (g *genomeNoMemory) HasY() bool {
    return g.y > 1
}

func (g *genomeNoMemory) Close() error {
    g.rsid = nil // GC map
    return g.m.Unmap()
}
