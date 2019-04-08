package game

import (
	"BountyWarServerG/game/model"
	"fmt"
	"sync"
)

type Chunk struct {
	X              int
	Y              int
	Num            uint16
	DronePool      []model.Drone
	DronePoolMutex sync.Mutex
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
			fmt.Println("remove from", obj.Param().CX, obj.Param().CY)
			fmt.Println("chunk size:", len(ch.DronePool))
			break
		}
	}
	DronePoolMutex.Unlock()
}

func (ch *Chunk) AddDrone(dr model.Drone) {
	DronePoolMutex.Lock()
	ch.DronePool = append(ch.DronePool, dr)
	DronePoolMutex.Unlock()
	fmt.Println("add to", dr.Param().CX, dr.Param().CY)
	fmt.Println("chunk size:", len(ch.DronePool))
}
