package tttandme

import (
    "unsafe"

    "github.com/coocood/freecache"
)

type genomeNoMemCached struct {
    timeout int
    g       genomeNoMemory
    cache   *freecache.Cache
}


func NewParserNoMemCached(size, timeoutSecs int) GenomeParser {
    return &genomeNoMemCached{
        cache: freecache.NewCache(size),
        timeout: timeoutSecs,
    }
}

func (g *genomeNoMemCached) Open(filename string) error {
    return g.g.Open(filename)
}

func (g *genomeNoMemCached) Parse() (Genome, error) {
    return g.g.Parse()
}

func (g *genomeNoMemCached) Close() error {
    return g.g.Close()
}

func (g *genomeNoMemCached) HasY() bool {
    return g.g.HasY()
}

func (g *genomeNoMemCached) Iter(f func(string) bool) bool {
    return g.g.Iter(f)
}

func (g *genomeNoMemCached) RSID(rsid string) *SNP {
    key := []byte(rsid)
    cached, err := g.cache.Get(key)

    // hit
    if err == nil {
        return (*SNP)(unsafe.Pointer(&cached[0]))
    }

    snp := g.g.RSID(rsid)
    if snp == nil {
        return nil
    }

    snpUnsafe := (*[10]byte)(unsafe.Pointer(snp))
    g.cache.Set(key, (*snpUnsafe)[:], g.timeout)

    return snp
}
