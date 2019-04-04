package game

import "BountyWarServerG/game/model"

type W_V_gramadina struct {
	vparam   *model.WeaponParameters
	rotation float64
	tip      int
}

func (t *W_V_gramadina) Update(delta int64) {
	//panic("implement me")
}

func (t *W_V_gramadina) Shoot() {
	//panic("implement me")
}

func (t *W_V_gramadina) Param() *model.WeaponParameters {
	return t.vparam
}
