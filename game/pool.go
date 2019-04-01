package game

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var PlayerPool map[uint32]*Player
var SessionPool map[*websocket.Conn]*Player
var PlayerPoolMutex sync.Mutex
var ObjectPool map[uint32]GameObject
var ObjectPoolMutex sync.Mutex

func init(){
	PlayerPool = make(map[uint32]*Player)
	ObjectPool = make(map[uint32]GameObject)
	SessionPool = make(map[*websocket.Conn]*Player)
}

func PlayerConnect(c *websocket.Conn) *Player {
	id := getNewID()
	wp := WeaponParameters{
		rotation:0,
	}
	w := TWeapon{
		vparam:&wp,
	}
	dp := DroneParameters{
		x:        100,
		y:        100,
		rotation: 0,
		iD:       id,
	}
	d := TDrone{
		parameters: &dp,
		tip:    2,
		weapon: &w,
	}
	p := Player{
		Session:      c,
		CurrentDrone: &d,
		KeySet:       [5]bool{false, false, false, false, false},
		answerPool:   bytes.NewBuffer(make([]byte, 0)),
		incomePool:   bytes.NewBuffer(make([]byte, 0)),
	}
	PlayerPoolMutex.Lock()
	PlayerPool[id] = &p
	SessionPool[c] = &p
	PlayerPoolMutex.Unlock()
	p.answerPoolMutex.Lock()
	p.answerPool.Write([]byte{0})
	p.answerPool.Write(idToBytes(id))
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
	log.Print("Player connect: ", len(PlayerPool))
	return &p
}

func IncomeMessage(ms []byte, c *websocket.Conn){
	if p, OK := SessionPool[c]; OK {
		p.incomePoolMutex.Lock()
		p.incomePool.Write(ms)
		p.incomePoolMutex.Unlock()
	}
}

func PlayerDisconnect(p *Player){

}