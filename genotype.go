package tttandme

import "fmt"

func GenotypeStr(genotype string) Genotype {
    genotp := nucleobases[genotype[0]]

    if len(genotype) > 1 {
        genotp |= nucleobases[genotype[1]] << 3
    }

    return genotp
}

func (g Genotype) String() string {
    if g == 0 {
        return "--"
    }
    if g0 := (g & 0x38) >> 3; g0 != 0 {
        return fmt.Sprintf("%c%c", baseStr[g & 0x7], baseStr[g0])
    }
    return string(baseStr[g & 0x7])
}

func (g Genotype) Complement() Genotype {
    first := g & 0x7
    second := (g & 0x38) >> 3
    return first.comp() | (second.comp() << 3)
}

func (g Genotype) ReverseComplement() Genotype {
    first := g & 0x7
    second := (g & 0x38) >> 3
    return (first.comp() << 3) | second.comp()
}

func (g Genotype) comp() Genotype {
    switch g {
    default:
        return g
    case BASE_A:
        return BASE_T
    case BASE_T:
        return BASE_A
    case BASE_C:
        return BASE_G
    case BASE_G:
        return BASE_C
    }
}

var nucleobases = map[byte]Genotype{
    '-': BASE_NONE,
    'A': BASE_A,
    'G': BASE_G,
    'C': BASE_C,
    'T': BASE_T,
    'D': BASE_D,
    'I': BASE_I,
}

var baseStr = map[Genotype]byte{
    BASE_NONE: '-',
    BASE_A: 'A',
    BASE_G: 'G',
    BASE_C: 'C',
    BASE_T: 'T',
    BASE_D: 'D',
    BASE_I: 'I',
}
