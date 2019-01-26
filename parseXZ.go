package tttandme

import (
    "os/exec"
    "io/ioutil"
    "strings"
    "errors"
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

    var sb strings.Builder

    xz := exec.Command("xz", "-c", "-d", filename)
    xz.Stdout = f
    xz.Stderr = &sb
    err = xz.Run()
    if err != nil {
        s := sb.String()
        return errors.New(s[:len(s)-1])
    }

    return pXZ.p.Open(f.Name())
}

func (pXZ *parserXZ) Parse() (Genome, error) {
    return pXZ.p.Parse()
}

func (pXZ *parserXZ) Close() error {
    return pXZ.p.Close()
}
