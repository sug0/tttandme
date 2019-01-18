package tttandme

// Represents a parser for a Human genome. Calling
// Parse() or Close() before Open() is undefined behavior.
// If one GenomeParser is initialized from another, it is
// expected of the topmost GenomeParser to Close() its
// child. Calling Parse() more than once is implementation
// defined.
type GenomeParser interface {
    Open(filename string) error
    Parse() (Genome, error)
    Close() error
}

// Represents a Human genome.
type Genome interface {
    HasY() bool
    RSID(rsid string) *SNP
    Iter(func(rsid string) bool) bool
}

// Represents a Human chromossome.
type Chromosome uint8

// Represents a position in the Human genome.
type Position uint64

// Represents a DNA base pair. Lowest 3 bits are
// the first nucleobase.
type Genotype uint8

// Represents a SNP, a location in the Human genome
// that is known to vary between individuals.
type SNP struct {
    Chromosome Chromosome
    Position   Position
    Genotype   Genotype
}

// Enumeration of all nucleobases.
const (
    BASE_NONE Genotype = iota
    BASE_A
    BASE_G
    BASE_C
    BASE_T
    BASE_D
    BASE_I
)

// Enumeration of all chromossomes.
const (
    CHR_01 Chromosome = iota + 1
    CHR_02
    CHR_03
    CHR_04
    CHR_05
    CHR_06
    CHR_07
    CHR_08
    CHR_09
    CHR_10
    CHR_11
    CHR_12
    CHR_13
    CHR_14
    CHR_15
    CHR_16
    CHR_17
    CHR_18
    CHR_19
    CHR_20
    CHR_21
    CHR_22
    CHR_X
    CHR_Y
    CHR_MT
)
