package stores

import (
	"encoding/binary"
	"strconv"
)

var (
	sizeOfUInt64                        = 8
	boltByteOrder      binary.ByteOrder = binary.LittleEndian
	boltkeyUsersBucket                  = getUInt64Bytes(0)
	boltkeyTokenBucket                  = getUInt64Bytes(1)
)

func getUInt64Bytes(v uint64) []byte {
	result := make([]byte, sizeOfUInt64)
	boltByteOrder.PutUint64(result, v)
	return result
}

func getStringFromUInt64Bytes(uBytes []byte) string {
	return strconv.FormatUint(boltByteOrder.Uint64(uBytes), 10)
}
