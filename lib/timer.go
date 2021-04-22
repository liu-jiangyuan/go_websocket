package lib

import (
	"log"
	"sync"
	"time"
)

type timer struct {
	id int64
	mutex sync.Mutex
	stopChan chan int64
	timerMap map[int64]*time.Timer

}


var Timer timer
func init() {
	Timer = timer{
		id:       0,
		mutex:    sync.Mutex{},
		stopChan: make(chan int64),
		timerMap: make(map[int64]*time.Timer),
	}
	go func() {
		for {
			select {
			case id := <- Timer.stopChan:
				Timer.Clear(id)
			}
		}
	}()
}


func (t *timer) newTimerId() int64 {
	t.mutex.Lock()
	t.id ++
	t.mutex.Unlock()
	return t.id
}


//清除一个定时器
func (t *timer) Clear(timerId int64) {
	t.timerMap[timerId].Stop()
	delete(t.timerMap,timerId)
}

//停止全部定时器
func (t *timer) ClearAll() {
	t.mutex.Lock()
	for timerId , _ := range t.timerMap{
		t.Clear(timerId)
	}
	t.mutex.Unlock()
}

//一次性定时器
func (t *timer) After (second time.Duration,action func() ) int64 {
	timerId := t.newTimerId()
	go func(id int64,) {
		t.timerMap[id] = time.NewTimer(second)
		for {
			select{
			case <-t.timerMap[id].C:
				action()
			}
		}
	}(timerId)
	return timerId
}

//周期定时器
func (t *timer)Tick(second time.Duration,action func()) int64 {
	timerId := t.newTimerId()
	go func(id int64) {
		restTimes := 0
		t.timerMap[id] = time.NewTimer(second)
		LOOP:
		for {
			select{
			case <- t.timerMap[id].C:
				if restTimes > 3 {
					t.Clear(id)
					break LOOP
				}
				t.timerMap[id].Reset(second)
				action()
				restTimes ++
			}
		}
	}(timerId)
	return timerId
}

func test() {
	id := Timer.After(time.Second * 3, func() {
		log.Printf("timer After func:%+v","1")
	})
	log.Printf("timer After id:%+v",id)
	id2 := Timer.After(time.Second * 5, func() {
		log.Printf("timer After func:%+v","id2")
	})
	log.Printf("timer Tick id:%+v",id2)

	id3 := Timer.Tick(time.Second * 7, func() {
		log.Printf("timer After func:%+v","id3")
	})
	log.Printf("timer Tick id:%+v",id3)



	Timer.Clear(id)
}
