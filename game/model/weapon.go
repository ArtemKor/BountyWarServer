package model

type Weapon interface {
	Param() *WeaponParameters
	Update(game Game)
	Shoot(game Game)
}

type WeaponParameters struct {
	Rotation float64
}
