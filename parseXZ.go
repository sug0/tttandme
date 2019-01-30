package tttandme

// Parses a xz compressed file; a temporary file will be used,
// created at the user's temporary file directory.
func NewParserXZ(p GenomeParser) GenomeParser {
    return NewParserCmd(p, "xz", "-c", "-d", "{}")
}
