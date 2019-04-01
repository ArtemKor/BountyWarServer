package game

import (
	"bytes"
	"github.com/gorilla/websocket"
	"sync"
)

const(
	KeyGo = 1
	KeyBack = 1 << 1
	KeyLeft = 1 << 2
	KeyRight = 1 << 3
	KeyFire = 1 << 4
)

var counterID uint32

var counterIDMutex sync.Mutex

func getNewID() uint32 {
	counterIDMutex.Lock()
	answer := counterID
	counterID++
	counterIDMutex.Unlock()
	return answer
}

type Player struct {
	Session         *websocket.Conn
	CurrentDrone    Drone
	KeySet          [5]bool
	answerPool      *bytes.Buffer
	answerPoolMutex sync.Mutex
	incomePool      *bytes.Buffer
	incomePoolMutex sync.Mutex
}

func (p *Player)onClose(code int, text string) error{
	return nil
}

func OnClose(code int, text string) error{
	return nil
}

func (p *Player) setKeys(kSet byte){
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

type GameObject interface {
	getID()uint32
	getDisplayProtocol()[]byte
	getPositionProtocol()[]byte
	update()
}

type Drone interface {
	param()*DroneParameters
}

type DroneParameters struct {
	x, y, rotation, mouse float64
	iD uint32
}


type Weapon interface {
	param() *WeaponParameters
	Update(delta int64)
	Shoot()
}

type WeaponParameters struct {
	rotation float64
}


type TWeapon struct{
	vparam *WeaponParameters
	rotation float64
	tip int
}

func (t *TWeapon) Update(delta int64) {
	//panic("implement me")
}

func (t *TWeapon) Shoot() {
	//panic("implement me")
}

func (t *TWeapon) param() *WeaponParameters{
	return t.vparam
}


type TDrone struct{
	parameters *DroneParameters
	tip int
	weapon Weapon
}

func (td *TDrone) update(){
	i := idToBytes(td.parameters.iD)
	b := td.getPositionProtocol()
	a := []byte{1,i[0],i[1],i[2],i[3]}
	a = append(a, b ...)
	broadcastMessage(a)
}

func (td *TDrone) param() *DroneParameters {
	return td.parameters
}

func (td *TDrone) getID() uint32 {
	return td.parameters.iD
}

func (td *TDrone) getDisplayProtocol() []byte {
	return []byte{2}
}

func (td *TDrone) getPositionProtocol() []byte {
	answer := make([]byte,10)
	x := coordToBytes(td.parameters.x)
	answer[0] = x[0]
	answer[1] = x[1]
	answer[2] = x[2]
	y := coordToBytes(td.parameters.y)
	answer[3] = y[0]
	answer[4] = y[1]
	answer[5] = y[2]
	answer[6] = 0
	answer[7] = 0
	a := radianToBytes(td.param().mouse)
	answer[8] = a[0]
	answer[9] = a[1]
	return answer
}
