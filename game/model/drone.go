package model

type Drone interface {
	Param() *DroneParameters
	GetID() uint32
	GetDisplayProtocol() []byte
	GetPositionProtocol() []byte
	Update()
}

type DroneParameters struct {
	X, Y, Rotation float64
	ID             uint32
	Health         int32
	Driver         *Player
}
