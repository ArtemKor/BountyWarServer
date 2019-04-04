package game

import (
	"BountyWarServerG/game/model"
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var SessionPool map[*websocket.Conn]*model.Player
var SessionPoolMutex sync.Mutex
var IncomePool map[*websocket.Conn]*model.Income
var IncomePoolMutex sync.Mutex
var DronePool map[uint32]model.Drone
var DronePoolMutex sync.Mutex
var ChunkPool map[uint16]*model.Chunk
var ChunkPoolMutex sync.Mutex

func init() {
	DronePool = make(map[uint32]model.Drone)
	SessionPool = make(map[*websocket.Conn]*model.Player)
	IncomePool = make(map[*websocket.Conn]*model.Income)
	ChunkPool = make(map[uint16]*model.Chunk)
	for x := 0; x < model.WorldSize+4; x++ {
		for y := 0; y < model.WorldSize+4; y++ {
			cnk := model.Chunk{
				X:   x,
				Y:   y,
				Num: uint16(x + y*256),
			}
			ChunkPool[cnk.Num] = &cnk
		}
	}
}

func PlayerConnect(c *websocket.Conn) *model.Player {
	id := getNewID()
	wp := model.WeaponParameters{
		Rotation: 0,
	}
	w := W_V_gramadina{
		vparam: &wp,
	}
	dp := model.DroneParameters{
		X:        100,
		Y:        100,
		Rotation: 0,
		ID:       id,
		Health:   100,
	}
	d := D_VI_geksa{
		parameters: &dp,
		tip:        2,
		weapon:     &w,
	}
	p := model.Player{
		Session:      c,
		CurrentDrone: &d,
		KeySet:       [5]bool{false, false, false, false, false},
		AnswerPool:   bytes.NewBuffer(make([]byte, 0)),
	}
	i := model.Income{
		Income: bytes.NewBuffer(make([]byte, 0)),
	}
	p.CurrentDrone.Param().Driver = &p
	IncomePoolMutex.Lock()
	IncomePool[c] = &i
	IncomePoolMutex.Unlock()
	SessionPoolMutex.Lock()
	SessionPool[c] = &p
	SessionPoolMutex.Unlock()
	p.AnswerPoolMutex.Lock()
	p.AnswerPool.Write([]byte{0})
	p.AnswerPool.Write(idToBytes(id))
	ma := make([]byte, 1025)
	ma[0] = 14
	ma[5] = 1
	ma[37] = 1
	ma[69] = 1
	p.AnswerPool.Write(ma)
	p.AnswerPoolMutex.Unlock()
	DronePoolMutex.Lock()
	DronePool[id] = &d
	DronePoolMutex.Unlock()
	//ms := []byte{0}
	//ms = append(ms, idToBytes(id) ...)
	//err := c.WriteMessage(0, ms)
	//if err != nil {
	//	log.Println("err:", err)
	//}
	log.Print("Player connect: ", len(SessionPool))
	return &p
}

func IncomeMessage(ms []byte, c *websocket.Conn) {
	if p, OK := IncomePool[c]; OK {
		p.IncomeMutex.Lock()
		p.Income.Write(ms)
		p.IncomeMutex.Unlock()
	}
}

func PlayerDisconnect(c *websocket.Conn) {
	SessionPoolMutex.Lock()
	p, OK := SessionPool[c]
	if OK {
		p.CurrentDrone.Param().Health = 0
		delete(SessionPool, c)
		log.Print("Player disconnect: ", len(SessionPool))
	}
	SessionPoolMutex.Unlock()
	_, OK = IncomePool[c]
	if OK {
		IncomePoolMutex.Lock()
		delete(IncomePool, c)
		IncomePoolMutex.Unlock()
	}
}
