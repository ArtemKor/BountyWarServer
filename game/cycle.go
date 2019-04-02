package game

import (
	"bytes"
	"github.com/gorilla/websocket"
	"time"
)

func Cycle() {
	current := time.Now()
	for {
		delta := current.UnixNano() - time.Now().UnixNano()
		current = time.Now()

		cin := make(map[*websocket.Conn]*bytes.Buffer)
		//IncomePoolMutex.Lock()
		for c, player := range IncomePool {
			if player.income.Len() > 0 {
				player.incomeMutex.Lock()
				buf := player.income
				player.income = bytes.NewBuffer(make([]byte, 0))
				player.incomeMutex.Unlock()
				cin[c] = buf
			}
		}
		//IncomePoolMutex.Unlock()

		for c, income := range cin {
			if income.Len() > 0 {
				incomeAnalise(income.Bytes(), c)
			}
		}

		for _, obj := range ObjectPool {
			obj.update()

		}

		for _, player := range SessionPool {
			if player.answerPool.Len() > 0 {
				player.answerPoolMutex.Lock()
				buf := player.answerPool
				player.answerPool = bytes.NewBuffer(make([]byte, 0))
				player.answerPoolMutex.Unlock()
				//log.Println("write :", buf.Len())
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
