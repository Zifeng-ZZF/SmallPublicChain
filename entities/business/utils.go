package business

import (
	"bytes"
	"encoding/binary"
	"log"
)

func getIntBytes(num int64) []byte {
	arr := new(bytes.Buffer)
	err := binary.Write(arr, binary.LittleEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return arr.Bytes()
}
