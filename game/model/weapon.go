package model

type Weapon interface {
	Param() *WeaponParameters
	Update(delta int64)
	Shoot()
}

type WeaponParameters struct {
	Rotation float64
}
