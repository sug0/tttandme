package main

import (
    "os"
    "fmt"
    "time"

    "github.com/sugoiuguu/tttandme"
)

func main() {
    dt := time.Now()

    p := tttandme.NewParserXZ(tttandme.NewParserNoMem())
    err := p.Open(os.Args[1])

    if err != nil {
        panic(err)
    }
    fmt.Fprintf(os.Stderr, "Decompressed in %s\n", time.Since(dt))
    defer p.Close()

    dt = time.Now()

    genome, err := p.Parse()
    if err != nil {
        panic(err)
    }
    fmt.Fprintf(os.Stderr, "Parsed in %s\n", time.Since(dt))

    var sex, smoke string

    if genome.HasY() {
        sex = "Male individual"
    } else {
        sex = "Female individual"
    }

    snpDep := genome.RSID("rs3750344")

    if snpDep == nil {
        fmt.Println("Unable to determine")
        return
    }
    dep := snpDep.Genotype.Complement()

    if dep == tttandme.GenotypeStr("AA") {
        smoke = "addicted"
    } else {
        smoke = "not addicted"
    }

    snpCancer := genome.RSID("rs1051730")

    if snpCancer == nil {
        fmt.Println("Unable to determine")
        return
    }
    geno := snpCancer.Genotype.Complement()

    switch {
    default:
        fmt.Println("Unable to determine")
    case geno == tttandme.GenotypeStr("CC"):
        fmt.Println(sex, "is likely", smoke, "and doesn't smoke much if a smoker")
    case geno == tttandme.GenotypeStr("CT"):
        fmt.Println(sex, "is likely", smoke, "and has 1.3x increased risk of lung cancer")
    case geno == tttandme.GenotypeStr("TT"):
        fmt.Println(sex, "is likely", smoke, "and has 1.8x increased risk of lung cancer")
    }
}
