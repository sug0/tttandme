package tttandme

import (
    "os/exec"
    "io/ioutil"
)

type parserXZ struct {
    p GenomeParser
}

func NewParserXZ(p GenomeParser) GenomeParser {
    return &parserXZ{p}
}

func (pXZ *parserXZ) Open(filename string) error {
    f, err := ioutil.TempFile("", "tttandme")
    if err != nil {
        return err
    }
    defer f.Close()

    xz := exec.Command("xz", "-c", "-d", filename)
    xz.Stdout = f
    err = xz.Run()
    if err != nil {
        return err
    }

    return pXZ.p.Open(f.Name())
}

func (pXZ *parserXZ) Parse() (Genome, error) {
    return pXZ.p.Parse()
}

func (pXZ *parserXZ) Close() error {
    return pXZ.p.Close()
}
