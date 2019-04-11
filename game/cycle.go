package game

import (
	"BountyWarServerG/game/model"
	"bytes"
	"math/rand"
	"time"
)

var game = MainGame{}

func initGame() {
	for c, cha := range ChunkPool {
		x := c % model.RealWorldSize
		y := c / model.RealWorldSize
		cha.Cells = make([]*ChunkCell, model.ChunkSize*model.ChunkSize)
		for i := 0; i < len(cha.Cells); i++ {
			cell := ChunkCell{}
			if x == 0 || x == model.RealWorldSize-1 || y == 0 || y == model.RealWorldSize-1 {
				cell.Tpe = 2
				cell.Grounded = false
			} else {
				cell.Tpe = 0
				cell.Grounded = true
			}
			cha.Cells[i] = &cell
		}
	}
	for i := 0; i < 10; i++ {
		x := rand.Intn(model.ChunkSize * model.WorldSize)
		y := rand.Intn(model.ChunkSize * model.WorldSize)
		Cx := x/model.ChunkSize + 1
		Cy := y/model.ChunkSize + 1
		cell := ChunkPool[Cx+model.RealWorldSize*Cy].Cells[(x%model.ChunkSize)+model.ChunkSize*(y%model.ChunkSize)]
		cell.Grounded = false
		cell.Tpe = 2
		curr := (y+model.ChunkSize)*model.WorldCellSize + x + model.ChunkSize
		prew := curr
		nx := 0
		ny := 0
		for n := 0; n < 6; n++ {
			new := curr
			for new == curr || new == prew {
				side := rand.Intn(4)
				if side == 0 {
					nx = 1
					ny = 0
				} else if side == 1 {
					nx = -1
					ny = 0
				} else if side == 2 {
					nx = 0
					ny = 1
				} else if side == 3 {
					nx = 0
					ny = -1
				}
				new = (y+ny+model.ChunkSize)*model.WorldCellSize + x + nx + model.ChunkSize
			}
			x += nx
			y += ny
			prew = curr
			curr = new
			Cx = (x + model.ChunkSize) / model.ChunkSize
			Cy = (y + model.ChunkSize) / model.ChunkSize
			//fmt.Println("x:", x, " y:", y)
			//fmt.Println("Cx:", Cx, " Cy:", Cy)
			cell = ChunkPool[Cx+model.RealWorldSize*Cy].Cells[((x+model.ChunkSize)%model.ChunkSize)+model.ChunkSize*((y+model.ChunkSize)%model.ChunkSize)]
			cell.Grounded = false
			cell.Tpe = 2
		}
	}
}

func Cycle() {
	initGame()
	current := time.Now()
	for {
		delta := current.UnixNano() - time.Now().UnixNano()
		current = time.Now()

		cin := make([]*model.ConnAnalize, 0, 16)
		//IncomePoolMutex.Lock()
		for c, player := range IncomePool {
			if player.Income.Len() > 0 {
				player.IncomeMutex.Lock()
				conn := model.ConnAnalize{
					Conn:   c,
					Buffer: player.Income,
				}
				player.Income = bytes.NewBuffer(make([]byte, 0))
				player.IncomeMutex.Unlock()
				cin = append(cin, &conn)
			}
		}
		//IncomePoolMutex.Unlock()
		cn := make(chan bool, model.CoreCount)

		if len(IncomePool) < model.CoreCount {
			execIncome(cin, cn)
			_ = <-cn
		} else {
			size := len(IncomePool) / model.CoreCount
			for i := 0; i < model.CoreCount; i++ {
				if i < model.CoreCount-1 {
					go execIncome(cin[size*i:size*(i+1)], cn)
				} else {
					go execIncome(cin[size*i:], cn)
				}
			}
			for i := 0; i < model.CoreCount; i++ {
				_ = <-cn
			}
		}

		if len(ChunkPool) < model.CoreCount {
			chunkDroneUpdate(ChunkPool, cn)
			_ = <-cn
		} else {
			size := len(IncomePool) / model.CoreCount
			for i := 0; i < model.CoreCount; i++ {
				if i < model.CoreCount-1 {
					go chunkDroneUpdate(ChunkPool[size*i:size*(i+1)], cn)
				} else {
					go chunkDroneUpdate(ChunkPool[size*i:], cn)
				}
			}
			for i := 0; i < model.CoreCount; i++ {
				_ = <-cn
			}
		}

		chunkDroneEvent(DroneChangeChunkPool, cn)
		_ = <-cn
		DroneChangeChunkPool = make([]*model.DroneChangeChunkEvent, 0)

		for _, player := range SessionPool {
			if player.AnswerPool.Len() > 0 {
				player.AnswerPoolMutex.Lock()
				buf := player.AnswerPool
				player.AnswerPool = bytes.NewBuffer(make([]byte, 0))
				player.AnswerPoolMutex.Unlock()
				_ = player.Session.WriteMessage(2, buf.Bytes())
			}
		}

		delta = current.UnixNano() - time.Now().UnixNano()
		sl := 25000000 - delta
		if sl > 0 {
			time.Sleep(time.Duration(sl))
		}
	}
}

func chunkDroneEvent(evs []*model.DroneChangeChunkEvent, compl chan (bool)) {
	for _, ev := range evs {
		DroneChangeChunk(ev)
	}
	compl <- true
}

func chunkDroneUpdate(cnks []*Chunk, compl chan (bool)) {
	for _, cnk := range cnks {
		cnk.UpdateDrones()
	}
	compl <- true
}

func execIncome(arr []*model.ConnAnalize, compl chan (bool)) {
	for _, conn := range arr {
		if conn.Buffer.Len() > 0 {
			IncomeAnalise(conn.Buffer.Bytes(), conn.Conn)
		}
	}
	compl <- true
}
