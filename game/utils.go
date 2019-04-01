package game

import (
	"encoding/binary"
)

func idToBytes(id uint32) []byte{
	answer := make([]byte,4)
	binary.BigEndian.PutUint32(answer, id)
	return answer
}

func bytesToID(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func radianToBytes(angle float64) []byte {
	answer := make([]byte,2)
	f := int(angle)
	s := int((angle - float64(f)) * 100)
	answer[0] = byte(f)
	answer[1] = byte(s)
	return answer
}

func bytesToRadian(angle []byte) float64 {
	return float64(angle[0]) + (float64(angle[1])/100.0)
}

func coordToBytes(coord float64) []byte {
	answer := make([]byte,3)
	f := int(coord) / 2048
	coord -= float64(f*2048)
	s := int(coord)/64
	coord -= float64(s*64)
	t := int(coord * 4)
	answer[0] = byte(f)
	answer[1] = byte(s)
	answer[2] = byte(t)
	return answer
}

func incomeAnalise(pocket []byte, p *Player){

	for i := 0; i < len(pocket); {
		switch pocket[i] {
		case 5:
			id := bytesToID(pocket[i+1: i+5])
			drone := ObjectPool[id].getDisplayProtocol()
			p.answerPoolMutex.Lock()
			p.answerPool.Write([]byte{2, pocket[i+1], pocket[i+2], pocket[i+3], pocket[i+4],})
			p.answerPool.Write(drone)
			p.answerPoolMutex.Unlock()
			i += 5
		case 6:
			p.setKeys(pocket[i+1])
			a := bytesToRadian(pocket[i+2:i+4])
			p.CurrentDrone.param().mouse = a
			i += 4
		default:
			i = len(pocket)
		}
	}
}

func broadcastMessage(b []byte){
	for _,player := range PlayerPool{
		player.answerPoolMutex.Lock()
		player.answerPool.Write(b)
		player.answerPoolMutex.Unlock()
	}
}
