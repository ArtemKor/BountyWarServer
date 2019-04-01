package game

import (
	"bytes"
	"time"
)

func Cycle(){
	current := time.Now()
	for {
		delta := current.UnixNano() - time.Now().UnixNano()
		current = time.Now()

		for _, player := range PlayerPool{
			if player.incomePool.Len() > 0 {
				player.incomePoolMutex.Lock()
				buf := player.incomePool
				player.incomePool = bytes.NewBuffer(make([]byte, 0))
				player.incomePoolMutex.Unlock()
				//log.Println("write :", buf.Len())
				incomeAnalise(buf.Bytes(), player)
			}
		}

		for _, obj := range ObjectPool{
			obj.update()

		}

		for _, player := range PlayerPool{
			if player.answerPool.Len() > 0 {
				player.answerPoolMutex.Lock()
				buf := player.answerPool
				player.answerPool = bytes.NewBuffer(make([]byte, 0))
				player.answerPoolMutex.Unlock()
				//log.Println("write :", buf.Len())
				player.Session.WriteMessage(2, buf.Bytes())
			}
		}


		delta = current.UnixNano() - time.Now().UnixNano()
		sl := 25000000 - delta
		if sl > 0 {
			time.Sleep(time.Duration(sl))
		}
	}
}