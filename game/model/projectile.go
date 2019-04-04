package model

type Projectile interface {
	Param() *ProjectileParameters
	GetID() uint32
	GetDisplayProtocol() []byte
	GetPositionProtocol() []byte
	Update()
}

type ProjectileParameters struct {
	X, Y, Rotation float64
	ID             uint32
	Health         int32
}
