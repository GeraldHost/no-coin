package nocoin

import (
	"math"
	"fmt"
)

const (
	varInt4BytePrefix = "FD"
	varInt8BytePrefix = "FE"
	varInt16BytePrefix = "FF"
)

func EncodeVarInt(n int) string {
	var prefix string
	var padding int

	str := fmt.Sprintf("%02X", n)
    nBytes := math.Ceil(float64(len(str))/float64(2))

    if nBytes > 1 && nBytes <= 2 {
    	padding = 4
    	prefix = "FD"
    } else if nBytes > 2 && nBytes <= 4 {
    	padding = 8
    	prefix = "FE"
    } else if nBytes > 4 {
    	padding = 16
    	prefix = "FF"
    }

	return fmt.Sprintf("%s%0*s", prefix, padding, str)
}