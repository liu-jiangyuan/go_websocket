package lib

import (
	"log"
	"runtime"
	"sync"
	"time"
)

type timer struct {
	id int64
	mutex sync.Mutex
	stopChan chan int64
	timerMap map[int64]*time.Timer
}

type TimerStruct struct {
	D time.Duration
	Args map[string]interface{}
	Action func(args map[string]interface{}) map[string]interface{}
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

//必须执行一次的定时器，无法销毁,带参数与不带参数两种
func (t *timer) MustOnceAfterWithFunArgs(s TimerStruct) {
	go time.AfterFunc(s.D, func() {
		s.Action(s.Args)
		//res := s.Action(s.Args)
		//log.Printf("MustOnceAfterWithFunArgs:%+v",res)
	})
}
func (t *timer) MustOnceAfter(second time.Duration,action func()) {
	go time.AfterFunc(second, action)
}

//一次性定时器,可根据timerId提前销毁，带参数与不带参数两种
func (t *timer) AfterWithFuncArgs (s TimerStruct) int64 {
	timerId := t.newTimerId()
	go time.AfterFunc(s.D, func() {
		s.Action(s.Args)
		//res := s.Action(s.Args)
		//log.Printf("AfterWithFuncArgs:%+v",res)
	})
	return timerId
}
func (t *timer) After (second time.Duration,action func()) int64 {
	timerId := t.newTimerId()
	t.timerMap[timerId] = time.AfterFunc(second, action)
	return timerId
}

//周期定时器,带参数与不带参数两种
func (t *timer) LoopWithFuncArgs (s TimerStruct,stop ...int) int64 {
	timerId := t.newTimerId()
	go func(id int64,s TimerStruct) {
		restTimes := 0
		t.timerMap[id] = time.NewTimer(s.D)
		for {
			select{
			case <- t.timerMap[id].C:
				//循环终止，方便调试
				if len(stop) > 1 && stop[0] == 1 && stop[1] <= restTimes {
					t.Clear(id)
					//终止当前协程
					runtime.Goexit()
				}
				restTimes ++
				t.timerMap[id].Reset(s.D)
				s.Action(s.Args)
				//res := s.Action(s.Args)
				//log.Printf("LoopWithFuncArgs:%+v",res)
			}
		}
	}(timerId,s)
	return timerId
}
func (t *timer) Loop (second time.Duration,action func(),stop ...int) int64 {
	timerId := t.newTimerId()

	go func(id int64) {
		restTimes := 0
		t.timerMap[id] = time.NewTimer(second)
		for {
			select{
			case <- t.timerMap[id].C:
				//循环终止，方便调试
				if len(stop) > 1 && stop[0] == 1 && stop[1] <= restTimes {
					t.Clear(id)
					//终止当前协程
					runtime.Goexit()
				}
				restTimes ++
				t.timerMap[id].Reset(second)
				action()
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

	id3 := Timer.Loop (time.Second * 3, func() {
		log.Printf("timer After func:%+v","id3")
	})
	log.Printf("timer Tick id:%+v",id3)
	Timer.MustOnceAfter(time.Second * 15 , func() {
		log.Printf("timer MustOnceAfter func:%+v","id4")
	})
	Timer.Clear(id)

	//id := Timer.AfterWithFuncArgs(TimerStruct{
	//	D:      7 * time.Second,
	//	Args: map[string]interface{}{"method":"AfterWithFuncArgsNoClear"},
	//	Action: controller.Index,
	//})
	//log.Printf("timer After id:%+v",id)
	//id2 := Timer.AfterWithFuncArgs(TimerStruct{
	//	D:      1 * time.Second,
	//	Args: map[string]interface{}{"method":"AfterWithFuncArgs"},
	//	Action: controller.Index,
	//})
	//log.Printf("timer Tick id:%+v",id2)
	//
	//id3 := Timer.LoopWithFuncArgs (TimerStruct{
	//	D:      3 * time.Second,
	//	Args: map[string]interface{}{"method":"LoopWithFuncArgs"},
	//	Action: controller.Index,
	//})
	//Timer.MustOnceAfterWithFunArgs(TimerStruct{
	//	D:      3 * time.Second,
	//	Args: map[string]interface{}{"method":"MustOnceAfterWithFunArgs"},
	//	Action: controller.Index,
	//})
	//log.Printf("timer Tick id:%+v",id3)
	//Timer.Clear(id)
}
