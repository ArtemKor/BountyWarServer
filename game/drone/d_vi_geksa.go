package drone

import (
	"BountyWarServerG/game/model"
	"math"
)

type D_VI_geksa struct {
	parameters *model.DroneParameters
	tip        int
	weapon     model.Weapon
	speed      float64
}

func (td *D_VI_geksa) Update() {
	if td.parameters.Health > 0 {

		if td.parameters.Driver.KeySet[0] {
			if td.speed == 0 {
				td.speed = 0.5
			}
			if td.speed < 4 {
				td.speed *= 1.024
			} else if td.speed > 4 {
				td.speed = 4
			}
		} else if td.speed != 0 {
			td.speed *= 0.966
			if td.speed < 0.5 && !td.parameters.Driver.KeySet[1] {
				td.speed = 0
			}
		}
		if td.parameters.Driver.KeySet[2] {
			td.parameters.Rotation -= 0.03
		}
		if td.parameters.Driver.KeySet[3] {
			td.parameters.Rotation += 0.03
		}

		if td.speed > 0 {
			td.parameters.X += td.speed * math.Sin(td.parameters.Rotation)
			td.parameters.Y -= td.speed * math.Cos(td.parameters.Rotation)
		}

		atan := math.Atan(float64(td.parameters.Driver.My)/float64(td.parameters.Driver.Mx)) + math.Pi/2.0
		if td.parameters.Driver.Mx < 0 {
			atan = atan + math.Pi
		}
		td.weapon.Param().Rotation = atan

		i := idToBytes(td.parameters.ID)
		b := td.GetPositionProtocol()
		a := []byte{1, i[0], i[1], i[2], i[3]}
		a = append(a, b...)
		broadcastMessage(a)
	} else {
		i := idToBytes(td.parameters.ID)
		a := []byte{4, i[0], i[1], i[2], i[3]}
		broadcastMessage(a)
	}
}

func (td *D_VI_geksa) Param() *model.DroneParameters {
	return td.parameters
}

func (td *D_VI_geksa) GetID() uint32 {
	return td.parameters.ID
}

func (td *D_VI_geksa) GetDisplayProtocol() []byte {
	return []byte{2}
}

func (td *D_VI_geksa) GetPositionProtocol() []byte {
	answer := make([]byte, 11)
	answer[0] = 6
	x := coordToBytes(td.parameters.X)
	answer[1] = x[0]
	answer[2] = x[1]
	answer[3] = x[2]
	y := coordToBytes(td.parameters.Y)
	answer[4] = y[0]
	answer[5] = y[1]
	answer[6] = y[2]
	b := radianToBytes(td.Param().Rotation)
	answer[7] = b[0]
	answer[8] = b[1]
	a := radianToBytes(td.weapon.Param().Rotation)
	answer[9] = a[0]
	answer[10] = a[1]
	return answer
}
