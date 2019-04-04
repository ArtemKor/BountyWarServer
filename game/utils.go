package game

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
)

func IncomeAnalise(pocket []byte, c *websocket.Conn) {
	p, OK := SessionPool[c]
	if OK {
		for i := 0; i < len(pocket); {
			switch pocket[i] {
			case 5:
				id := bytesToID(pocket[i+1 : i+5])
				drone := DronePool[id].getDisplayProtocol()
				p.answerPoolMutex.Lock()
				p.answerPool.Write([]byte{2, pocket[i+1], pocket[i+2], pocket[i+3], pocket[i+4]})
				p.answerPool.Write(drone)
				p.answerPoolMutex.Unlock()
				i += 5
			case 6:
				p.setKeys(pocket[i+1])
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

func broadcastMessage(b []byte) {
	for _, player := range SessionPool {
		player.answerPoolMutex.Lock()
		player.answerPool.Write(b)
		player.answerPoolMutex.Unlock()
	}
}
