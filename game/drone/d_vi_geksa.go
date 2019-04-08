package drone

import (
	"BountyWarServerG/game/model"
	"math"
)

type D_VI_geksa struct {
	Parameters *model.DroneParameters
	Tip        int
	Weapon     model.Weapon
	speed      float64
}

func (td *D_VI_geksa) Update(game model.Game) {
	if td.Parameters.Health > 0 {

		if td.Parameters.Driver.KeySet[0] {
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
			if td.speed < 0.5 && !td.Parameters.Driver.KeySet[1] {
				td.speed = 0
			}
		}
		if td.Parameters.Driver.KeySet[2] {
			td.Parameters.Rotation -= 0.03
		}
		if td.Parameters.Driver.KeySet[3] {
			td.Parameters.Rotation += 0.03
		}

		if td.speed > 0 {
			td.Parameters.X += td.speed * math.Sin(td.Parameters.Rotation)
			td.Parameters.Y -= td.speed * math.Cos(td.Parameters.Rotation)
			nCX := uint16(td.Parameters.X / model.ChunkMulty)
			nCY := uint16(td.Parameters.Y / model.ChunkMulty)
			if td.Parameters.CX != nCX || td.Parameters.CY != nCY {
				game.AddDroneChangeChunkEvent(nCX, nCY, td)
			}
		}

		atan := math.Atan(float64(td.Parameters.Driver.My)/float64(td.Parameters.Driver.Mx)) + math.Pi/2.0
		if td.Parameters.Driver.Mx < 0 {
			atan = atan + math.Pi
		}
		td.Weapon.Param().Rotation = atan

		i := model.IdToBytes(td.Parameters.ID)
		b := td.GetPositionProtocol()
		a := []byte{1, i[0], i[1], i[2], i[3]}
		a = append(a, b...)
		game.BroadcastMessage(a)
	} else {
		game.RemoveDrone(td)
		i := model.IdToBytes(td.Parameters.ID)
		a := []byte{4, i[0], i[1], i[2], i[3]}
		game.BroadcastMessage(a)
	}
}

func (td *D_VI_geksa) Param() *model.DroneParameters {
	return td.Parameters
}

func (td *D_VI_geksa) GetID() uint32 {
	return td.Parameters.ID
}

func (td *D_VI_geksa) GetDisplayProtocol() []byte {
	return []byte{2}
}

func (td *D_VI_geksa) GetPositionProtocol() []byte {
	answer := make([]byte, 11)
	answer[0] = 6
	x := model.CoordToBytes(td.Parameters.X)
	answer[1] = x[0]
	answer[2] = x[1]
	answer[3] = x[2]
	y := model.CoordToBytes(td.Parameters.Y)
	answer[4] = y[0]
	answer[5] = y[1]
	answer[6] = y[2]
	b := model.RadianToBytes(td.Param().Rotation)
	answer[7] = b[0]
	answer[8] = b[1]
	a := model.RadianToBytes(td.Weapon.Param().Rotation)
	answer[9] = a[0]
	answer[10] = a[1]
	return answer
}
