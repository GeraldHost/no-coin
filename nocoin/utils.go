package nocoin

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
    "crypto/sha256"
)

var (
	varInt2BytePrefix  []byte = []byte("")
	varInt4BytePrefix  []byte = []byte("FD")
	varInt8BytePrefix  []byte = []byte("FE")
	varInt16BytePrefix []byte = []byte("FF")
)

// Decode varint
// returns number and number of bits read
var varIntPrefixes map[string]int = map[string]int{
	string(varInt2BytePrefix):  2,
	string(varInt4BytePrefix):  4,
	string(varInt8BytePrefix):  8,
	string(varInt16BytePrefix): 16}

func VarIntFromReader(r *bytes.Buffer) int64 {
    maybePrefix := r.Next(2)
    var num int64
    if nBytes, ok := varIntPrefixes[string(maybePrefix)]; ok {
        num, _ = BytesToInt(r.Next(nBytes))
    } else {
        num, _ = BytesToInt(maybePrefix)
    }
    return num
}

func DecodeVarInt(s string) (int64, int, error) {
	b := []byte(s)
	startingChar := 0
	prefix := varInt2BytePrefix
	if bytes.HasPrefix(b, varInt4BytePrefix) || bytes.HasPrefix(b, varInt8BytePrefix) || bytes.HasPrefix(b, varInt16BytePrefix) {
		startingChar = 2
		prefix = b[0:2]
	}
	intSize := varIntPrefixes[string(prefix)]
	hex := b[startingChar : startingChar+intSize]
	n, err := strconv.ParseInt(string(hex), 16, 64)
	if err != nil {
		return 0, 0, errors.New("invalid varint")
	}
	return n, intSize + startingChar, nil
}

// TODO: convert this to return []byte
func EncodeVarInt(n int) string {
	var prefix []byte
	str := fmt.Sprintf("%02X", n)
	nBytes := math.Ceil(float64(len(str)) / float64(2))
	if nBytes > 1 && nBytes <= 2 {
		prefix = varInt4BytePrefix
	} else if nBytes > 2 && nBytes <= 4 {
		prefix = varInt8BytePrefix
	} else if nBytes > 4 {
		prefix = varInt16BytePrefix
	}
	padding := varIntPrefixes[string(prefix)]
	return fmt.Sprintf("%s%0*s", prefix, padding, str)
}

// convery bytes to int
func BytesToInt(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 16, 64)
}

func Sha256(s string) string {
    hash := sha256.Sum256([]byte(s))
    return fmt.Sprintf("%x", hash)
}
