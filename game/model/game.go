package model

type Game interface {
	RemoveDrone(dr Drone)
	AddDroneChangeChunkEvent(nx, ny uint16, dr Drone)
	BroadcastMessage(b []byte)
}
