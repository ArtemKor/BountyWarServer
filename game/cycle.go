package game

import (
	"BountyWarServerG/game/model"
	"bytes"
	"time"
)

var game = MainGame{}

func Cycle() {
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
