package weapon

import "BountyWarServerG/game/model"

type W_V_gramadina struct {
	Vparam   *model.WeaponParameters
	Rotation float64
	Tip      int
}

func (t *W_V_gramadina) Update(game model.Game) {
	//panic("implement me")
}

func (t *W_V_gramadina) Shoot(game model.Game) {
	//panic("implement me")
}

func (t *W_V_gramadina) Param() *model.WeaponParameters {
	return t.Vparam
}
