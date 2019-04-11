package game

import (
	"BountyWarServerG/game/drone"
	"BountyWarServerG/game/model"
	"BountyWarServerG/game/weapon"
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var SessionPool map[*websocket.Conn]*model.Player
var PlayerPool map[*model.Player]*websocket.Conn
var SessionPoolMutex sync.Mutex
var IncomePool map[*websocket.Conn]*model.Income
var IncomePoolMutex sync.Mutex
var DronePool map[uint32]model.Drone
var DronePoolMutex sync.Mutex
var ChunkPool []*Chunk
var ChunkPoolMutex sync.Mutex
var DroneChangeChunkPool []*model.DroneChangeChunkEvent
var DroneChangeChunkPoolMutex sync.Mutex

func init() {
	DronePool = make(map[uint32]model.Drone)
	SessionPool = make(map[*websocket.Conn]*model.Player)
	PlayerPool = make(map[*model.Player]*websocket.Conn)
	IncomePool = make(map[*websocket.Conn]*model.Income)
	ChunkPool = make([]*Chunk, model.RealWorldSize*model.RealWorldSize)
	for x := 0; x < model.RealWorldSize; x++ {
		for y := 0; y < model.RealWorldSize; y++ {
			cnk := Chunk{
				X:         x,
				Y:         y,
				Num:       uint16(x + y*model.RealWorldSize),
				DronePool: make([]model.Drone, 0),
			}
			ChunkPool[cnk.Num] = &cnk
		}
	}
}

func PlayerConnect(c *websocket.Conn) *model.Player {

	p := model.Player{
		Session:      c,
		CurrentDrone: nil,
		KeySet:       [5]bool{false, false, false, false, false},
		AnswerPool:   bytes.NewBuffer(make([]byte, 0)),
		Chunks:       make([]bool, model.RealWorldSize*model.RealWorldSize),
	}
	i := model.Income{
		Income: bytes.NewBuffer(make([]byte, 0)),
	}
	IncomePoolMutex.Lock()
	IncomePool[c] = &i
	IncomePoolMutex.Unlock()
	SessionPoolMutex.Lock()
	SessionPool[c] = &p
	PlayerPool[&p] = c
	SessionPoolMutex.Unlock()
	//ma := make([]byte, 1025)
	//ma[0] = 14
	//ma[5] = 1
	//ma[37] = 1
	//ma[69] = 1
	//p.AnswerPoolMutex.Lock()
	//p.AnswerPool.Write(ma)
	//p.AnswerPoolMutex.Unlock()
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
		if p.CurrentDrone != nil {
			p.CurrentDrone.Param().Health = 0
		}

		delete(PlayerPool, SessionPool[c])
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

func CreateDroneForPlayer(c *websocket.Conn, p *model.Player) {
	id := model.GetNewID()
	wp := model.WeaponParameters{
		Rotation: 0,
	}
	w := weapon.W_V_gramadina{
		Vparam: &wp,
	}
	dp := model.DroneParameters{
		X:        50 + model.ChunkMulty,
		Y:        50 + model.ChunkMulty,
		Rotation: 0,
		ID:       id,
		Health:   100,
	}
	d := drone.D_VI_geksa{
		Parameters: &dp,
		Tip:        2,
		Weapon:     &w,
	}
	p.CurrentDrone = &d
	p.CurrentDrone.Param().Driver = p

	p.AnswerPoolMutex.Lock()
	p.AnswerPool.Write([]byte{0})
	p.AnswerPool.Write(model.IdToBytes(id))
	p.AnswerPoolMutex.Unlock()

	DronePoolMutex.Lock()
	DronePool[id] = &d
	DronePoolMutex.Unlock()

	//fmt.Printlln("New drone")
	game.AddDroneChangeChunkEvent(uint16(dp.X/model.ChunkMulty), uint16(dp.Y/model.ChunkMulty), &d)
	//ms := []byte{0}
	//ms = append(ms, model.IdToBytes(id) ...)
	//err := c.WriteMessage(1, ms)
	//if err != nil {
	//	log.Println("err:", err)
	//}
}
