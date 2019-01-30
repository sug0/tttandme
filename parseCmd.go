package tttandme

import (
    "os/exec"
    "io/ioutil"
    "strings"
    "errors"
)

type parserCmd struct {
    p    GenomeParser
    name string
    args []string
}

// Returns a new GenomeParser that will run a command on the input filename,
// creating a temporary file with the output where the GenomeParser p will run.
// The string "{}" will be replaced by the filename argument.
func NewParserCmd(p GenomeParser, name string, arg ...string) GenomeParser {
    return &parserCmd{p, name, arg}
}

func (p *parserCmd) Open(filename string) error {
    f, err := ioutil.TempFile("", "tttandme")
    if err != nil {
        return err
    }
    defer f.Close()

    var sb strings.Builder

    xz := exec.Command(p.name, replaceFilename(filename, p.args)...)
    p.name = ""
    p.args = nil

    xz.Stdout = f
    xz.Stderr = &sb
    err = xz.Run()

    if err != nil {
        return errors.New(sb.String())
    }

    return p.p.Open(f.Name())
}

func (p *parserCmd) Parse() (Genome, error) {
    return p.p.Parse()
}

func (p *parserCmd) Close() error {
    return p.p.Close()
}

func replaceFilename(filename string, args []string) []string {
    for i := range args {
        if args[i] == "{}" {
            args[i] = filename
        }
    }
    return args
}
