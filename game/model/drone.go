package model

type Drone interface {
	Param() *DroneParameters
	GetID() uint32
	GetDisplayProtocol() []byte
	GetPositionProtocol() []byte
	Update(game Game)
}

type DroneParameters struct {
	X, Y, Rotation float64
	ID             uint32
	Health         int32
	CX, CY         uint16
	Driver         *Player
}
