package model

import (
	"bytes"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	KeyGo      = 1
	KeyBack    = 1 << 1
	KeyLeft    = 1 << 2
	KeyRight   = 1 << 3
	KeyFire    = 1 << 4
	CellSize   = 64
	ChunkSize  = 32
	WorldSize  = 3
	CellMulty  = 256.0 / CellSize
	ChunkMulty = CellSize * ChunkSize
	CoreCount  = 7
)

type GameObject interface {
	GetID() uint32
	GetDisplayProtocol() []byte
	GetPositionProtocol() []byte
	Update()
}

type ConnAnalize struct {
	Conn   *websocket.Conn
	Buffer *bytes.Buffer
}

type Player struct {
	Session         *websocket.Conn
	CurrentDrone    Drone
	KeySet          [5]bool
	Mx              int
	My              int
	AnswerPool      *bytes.Buffer
	AnswerPoolMutex sync.Mutex
}

type Income struct {
	Income      *bytes.Buffer
	IncomeMutex sync.Mutex
}

func (p *Player) SetKeys(kSet byte) {
	if kSet&KeyGo != 0 {
		p.KeySet[0] = true
	} else {
		p.KeySet[0] = false
	}
	if kSet&KeyBack != 0 {
		p.KeySet[1] = true
	} else {
		p.KeySet[1] = false
	}
	if kSet&KeyLeft != 0 {
		p.KeySet[2] = true
	} else {
		p.KeySet[2] = false
	}
	if kSet&KeyRight != 0 {
		p.KeySet[3] = true
	} else {
		p.KeySet[3] = false
	}
	if kSet&KeyFire != 0 {
		p.KeySet[4] = true
	} else {
		p.KeySet[4] = false
	}
}
