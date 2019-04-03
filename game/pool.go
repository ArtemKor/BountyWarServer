package game

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var SessionPool map[*websocket.Conn]*Player
var SessionPoolMutex sync.Mutex
var IncomePool map[*websocket.Conn]*Income
var IncomePoolMutex sync.Mutex
var ObjectPool map[uint32]GameObject
var ObjectPoolMutex sync.Mutex

func init() {
	ObjectPool = make(map[uint32]GameObject)
	SessionPool = make(map[*websocket.Conn]*Player)
	IncomePool = make(map[*websocket.Conn]*Income)
}

func PlayerConnect(c *websocket.Conn) *Player {
	id := getNewID()
	wp := WeaponParameters{
		rotation: 0,
	}
	w := TWeapon{
		vparam: &wp,
	}
	dp := DroneParameters{
		x:        100,
		y:        100,
		rotation: 0,
		iD:       id,
		health:   100,
	}
	d := TDrone{
		parameters: &dp,
		tip:        2,
		weapon:     &w,
	}
	p := Player{
		Session:      c,
		CurrentDrone: &d,
		KeySet:       [5]bool{false, false, false, false, false},
		answerPool:   bytes.NewBuffer(make([]byte, 0)),
	}
	i := Income{
		income: bytes.NewBuffer(make([]byte, 0)),
	}
	p.CurrentDrone.param().driver = &p
	IncomePoolMutex.Lock()
	IncomePool[c] = &i
	IncomePoolMutex.Unlock()
	SessionPoolMutex.Lock()
	SessionPool[c] = &p
	SessionPoolMutex.Unlock()
	p.answerPoolMutex.Lock()
	p.answerPool.Write([]byte{0})
	p.answerPool.Write(idToBytes(id))
	ma := make([]byte, 1025)
	ma[0] = 14
	ma[5] = 1
	ma[37] = 1
	ma[69] = 1
	p.answerPool.Write(ma)
	p.answerPoolMutex.Unlock()
	ObjectPoolMutex.Lock()
	ObjectPool[id] = &d
	ObjectPoolMutex.Unlock()
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
		p.incomeMutex.Lock()
		p.income.Write(ms)
		p.incomeMutex.Unlock()
	}
}

func PlayerDisconnect(c *websocket.Conn) {
	SessionPoolMutex.Lock()
	p, OK := SessionPool[c]
	if OK {
		p.CurrentDrone.param().health = 0
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
