package tttandme

import "fmt"

func (chrxm Chromosome) String() string {
    switch chrxm {
    default:
        return fmt.Sprintf("%d", chrxm)
    case CHR_X:
        return "X"
    case CHR_Y:
        return "Y"
    case CHR_MT:
        return "MT"
    }
}
