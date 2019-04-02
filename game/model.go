package game

import (
	"bytes"
	"github.com/gorilla/websocket"
	"math"
	"sync"
)

const (
	KeyGo    = 1
	KeyBack  = 1 << 1
	KeyLeft  = 1 << 2
	KeyRight = 1 << 3
	KeyFire  = 1 << 4
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
}

type Income struct {
	income      *bytes.Buffer
	incomeMutex sync.Mutex
}

func (p *Player) setKeys(kSet byte) {
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
	getID() uint32
	getDisplayProtocol() []byte
	getPositionProtocol() []byte
	update()
}

type Drone interface {
	param() *DroneParameters
}

type DroneParameters struct {
	x, y, rotation, mouse float64
	iD                    uint32
	health                int32
	driver                *Player
}

type Weapon interface {
	param() *WeaponParameters
	Update(delta int64)
	Shoot()
}

type WeaponParameters struct {
	rotation float64
}

type TWeapon struct {
	vparam   *WeaponParameters
	rotation float64
	tip      int
}

func (t *TWeapon) Update(delta int64) {
	//panic("implement me")
}

func (t *TWeapon) Shoot() {
	//panic("implement me")
}

func (t *TWeapon) param() *WeaponParameters {
	return t.vparam
}

type TDrone struct {
	parameters *DroneParameters
	tip        int
	weapon     Weapon
	speed      float64
}

func (td *TDrone) update() {
	if td.parameters.health > 0 {

		if td.parameters.driver.KeySet[0] {
			if td.speed == 0 {
				td.speed = 0.5
			}
			if td.speed < 4 {
				td.speed *= 1.024
			} else if td.speed > 4 {
				td.speed = 4
			}
		} else if td.speed != 0 {
			td.speed *= 0.966
			if td.speed < 0.5 && !td.parameters.driver.KeySet[1] {
				td.speed = 0
			}
		}
		if td.parameters.driver.KeySet[2] {
			td.parameters.rotation -= 0.03
		}
		if td.parameters.driver.KeySet[3] {
			td.parameters.rotation += 0.03
		}

		if td.speed > 0 {
			td.parameters.x += td.speed * math.Sin(td.parameters.rotation)
			td.parameters.y -= td.speed * math.Cos(td.parameters.rotation)
		}

		i := idToBytes(td.parameters.iD)
		b := td.getPositionProtocol()
		a := []byte{1, i[0], i[1], i[2], i[3]}
		a = append(a, b...)
		broadcastMessage(a)
	} else {
		i := idToBytes(td.parameters.iD)
		a := []byte{4, i[0], i[1], i[2], i[3]}
		broadcastMessage(a)
	}
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
	answer := make([]byte, 10)
	x := coordToBytes(td.parameters.x)
	answer[0] = x[0]
	answer[1] = x[1]
	answer[2] = x[2]
	y := coordToBytes(td.parameters.y)
	answer[3] = y[0]
	answer[4] = y[1]
	answer[5] = y[2]
	b := radianToBytes(td.param().rotation)
	answer[6] = b[0]
	answer[7] = b[1]
	a := radianToBytes(td.param().mouse)
	answer[8] = a[0]
	answer[9] = a[1]
	return answer
}
