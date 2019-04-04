package model

import (
	"encoding/binary"
	"sync"
)

var counterID uint32

var counterIDMutex sync.Mutex

func GetNewID() uint32 {
	counterIDMutex.Lock()
	answer := counterID
	counterID++
	counterIDMutex.Unlock()
	return answer
}

func IdToBytes(id uint32) []byte {
	answer := make([]byte, 4)
	binary.BigEndian.PutUint32(answer, id)
	return answer
}

func BytesToID(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func RadianToBytes(angle float64) []byte {
	answer := make([]byte, 2)
	f := int(angle)
	s := int((angle - float64(f)) * 100)
	answer[0] = byte(f)
	answer[1] = byte(s)
	return answer
}

func BytesToRadian(angle []byte) float64 {
	return float64(angle[0]) + (float64(angle[1]) / 100.0)
}

func CoordToBytes(coord float64) []byte {
	answer := make([]byte, 3)
	f := int(coord) / ChunkMulty
	//coord -= float64(f * 2048)
	s := (int(coord) % ChunkMulty) / CellSize
	//coord -= float64(s * 64)
	t := int(coord*CellMulty) % CellSize
	answer[0] = byte(f)
	answer[1] = byte(s)
	answer[2] = byte(t)
	return answer
}
