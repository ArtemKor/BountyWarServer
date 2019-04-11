package game

import (
	"BountyWarServerG/game/model"
	"sync"
)

type Chunk struct {
	X              int
	Y              int
	Num            uint16
	DronePool      []model.Drone
	DronePoolMutex sync.Mutex
	Cells          []*ChunkCell
}

type ChunkCell struct {
	Tpe      byte
	Grounded bool
}

func (ch *Chunk) GenerateGroundPocket() []byte {
	ansver := make([]byte, model.ChunkSize*model.ChunkSize+3)
	ansver[0] = 14
	ansver[1] = byte(ch.X)
	ansver[2] = byte(ch.Y)
	for i, b := range ch.Cells {
		ansver[i+3] = b.Tpe
	}
	return ansver
}

func (ch *Chunk) UpdateDrones() {
	for _, obj := range ch.DronePool {
		obj.Update(&game)
	}
}

func (ch *Chunk) RemoveDrone(dr model.Drone) {
	DronePoolMutex.Lock()
	for i, obj := range ch.DronePool {
		if obj == dr {
			ch.DronePool = append(ch.DronePool[:i], ch.DronePool[i+1:]...)
			//fmt.Println("remove from", obj.Param().CX, obj.Param().CY)
			//fmt.Println("chunk size:", len(ch.DronePool))
			break
		}
	}
	DronePoolMutex.Unlock()
}

func (ch *Chunk) AddDrone(dr model.Drone) {
	DronePoolMutex.Lock()
	ch.DronePool = append(ch.DronePool, dr)
	DronePoolMutex.Unlock()
	//fmt.Println("add to", dr.Param().CX, dr.Param().CY)
	//fmt.Println("chunk size:", len(ch.DronePool))
	driver := dr.Param().Driver
	num := ch.X + ch.Y*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ch.GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X - 1 + (ch.Y-1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X + (ch.Y-1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X + 1 + (ch.Y-1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X - 1 + (ch.Y)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X + 1 + (ch.Y)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X - 1 + (ch.Y+1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X + (ch.Y+1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
	num = ch.X + 1 + (ch.Y+1)*model.RealWorldSize
	if !driver.Chunks[num] {
		driver.AnswerPool.Write(ChunkPool[num].GenerateGroundPocket())
		driver.Chunks[num] = true
	}
}
