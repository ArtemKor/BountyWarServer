package game

import (
	"BountyWarServerG/game/model"
	"encoding/binary"
	"github.com/gorilla/websocket"
)

type MainGame struct{}

func (mg *MainGame) RemoveDrone(dr model.Drone) {
	ChunkPool[dr.Param().CX+dr.Param().CY*model.RealWorldSize].RemoveDrone(dr)
}

func DroneChangeChunk(ev *model.DroneChangeChunkEvent) {
	if ev.Cx == 0 || ev.Cx == model.RealWorldSize-1 || ev.Cy == 0 || ev.Cy == model.RealWorldSize-1 {
		ev.Drone.Param().Health = 0
	} else {
		pr := ev.Drone.Param()
		cnum := pr.CY*model.RealWorldSize + pr.CX
		nnum := ev.Cy*model.RealWorldSize + ev.Cx
		ChunkPool[cnum].RemoveDrone(ev.Drone)
		pr.CX = ev.Cx
		pr.CY = ev.Cy
		ChunkPool[nnum].AddDrone(ev.Drone)
	}
}

func (mg *MainGame) AddDroneChangeChunkEvent(nx, ny uint16, dr model.Drone) {
	pr := model.DroneChangeChunkEvent{
		Drone: dr,
		Cx:    nx,
		Cy:    ny,
	}
	DroneChangeChunkPoolMutex.Lock()
	DroneChangeChunkPool = append(DroneChangeChunkPool, &pr)
	DroneChangeChunkPoolMutex.Unlock()
}

func (mg *MainGame) BroadcastMessage(b []byte) {
	for _, player := range SessionPool {
		player.AnswerPoolMutex.Lock()
		player.AnswerPool.Write(b)
		player.AnswerPoolMutex.Unlock()
	}
}

func IncomeAnalise(pocket []byte, c *websocket.Conn) {
	p, OK := SessionPool[c]
	if OK {
		for i := 0; i < len(pocket); {
			switch pocket[i] {
			case 1:
				p := SessionPool[c]
				if p.CurrentDrone == nil {
					_ = pocket[i+1]
					CreateDroneForPlayer(c, SessionPool[c])
				}
				i += 2
			case 5:
				id := model.BytesToID(pocket[i+1 : i+5])
				drone := DronePool[id].GetDisplayProtocol()
				p.AnswerPoolMutex.Lock()
				p.AnswerPool.Write([]byte{2, pocket[i+1], pocket[i+2], pocket[i+3], pocket[i+4]})
				p.AnswerPool.Write(drone)
				p.AnswerPoolMutex.Unlock()
				i += 5
			case 6:
				p.SetKeys(pocket[i+1])
				p.Mx = int(binary.BigEndian.Uint16(pocket[i+2:i+4])) - 32768
				p.My = int(binary.BigEndian.Uint16(pocket[i+4:i+6])) - 32768
				//a := bytesToRadian(pocket[i+2 : i+4])
				//p.CurrentDrone.param().mouse = a
				i += 6
			default:
				i = len(pocket)
			}
		}
	}
}
