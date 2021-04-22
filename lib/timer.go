package lib

import (
	"github.com/liu-jiangyuan/go_websocket/controller"
	"log"
	"runtime"
	"sync"
	"time"
)

type timer struct {
	id int64
	mutex sync.Mutex
	stopChan chan int64
	timerMap map[int64]TimerStruct
}

type TimerStruct struct {
	D time.Duration
	t *time.Timer
	Args map[string]interface{}
	Action func(args map[string]interface{}) map[string]interface{}
}


var Timer timer
func init() {
	Timer = timer{
		id:       0,
		mutex:    sync.Mutex{},
		stopChan: make(chan int64),
		timerMap: make(map[int64]TimerStruct),
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
	if _, ok := t.timerMap[timerId];ok {
		t.timerMap[timerId].t.Stop()
		delete(t.timerMap,timerId)
	}
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
func (t *timer) MustOnceAfter (s TimerStruct) {
	time.AfterFunc(s.D, func() {
		//s.Action(s.Args)
		res := s.Action(s.Args)
		log.Printf("MustOnceAfter:%+v",res)
	})
}

//一次性定时器,可根据timerId提前销毁，带参数与不带参数两种

func (t *timer) After (s TimerStruct) int64 {
	timerId := t.newTimerId()
	t.timerMap[timerId] = TimerStruct{
		D:      s.D,
		t:      time.AfterFunc(s.D, func() {
			//s.Action(s.Args)
			res := s.Action(s.Args)
			log.Printf("After:%+v",res)}),
		Args:   s.Args,
		Action: s.Action,
	}
	return timerId
}

//周期定时器,带参数与不带参数两种
func (t *timer) Loop (s TimerStruct,stop ...int) int64 {
	timerId := t.newTimerId()
	t.timerMap[timerId] = TimerStruct{
		D:      s.D,
		t:      time.NewTimer(s.D),
		Args:   s.Args,
		Action: s.Action,
	}

	go func(id int64,s TimerStruct) {
		defer func() {
			if err := recover(); err != nil {
				s.t.Stop()
				t.Clear(id)
			}
		}()
		restTimes := 0
		for {
			select{
			case <- t.timerMap[id].t.C:
				//循环终止，方便调试
				if len(stop) > 1 && stop[0] == 1 && stop[1] <= restTimes {
					t.Clear(id)
					//终止当前协程
					runtime.Goexit()
				}
				restTimes ++
				t.timerMap[id].t.Reset(s.D)
				//s.Action(s.Args)
				res := s.Action(s.Args)
				log.Printf("Loop:%+v",res)
			}
		}
	}(timerId,t.timerMap[timerId])

	return timerId
}

func test() {
	id := Timer.After(TimerStruct{
		D:      7 * time.Second,
		Args: map[string]interface{}{"method":"AfterWithFuncArgsNoClear"},
		Action: controller.Index,
	})
	log.Printf("Timer After id:%+v",id)
	id2 := Timer.After(TimerStruct{
		D:      9 * time.Second,
		Args: map[string]interface{}{"method":"AfterWithFuncArgs"},
		Action: controller.Index,
	})
	log.Printf("Timer Tick id:%+v",id2)

	id3 := Timer.Loop (TimerStruct{
		D:      1 * time.Second,
		Args: map[string]interface{}{"method":"LoopWithFuncArgs"},
		Action: controller.Index,
	},1,5)
	log.Printf("Timer Tick id:%+v",id3)

	Timer.MustOnceAfter(TimerStruct{
		D:      15 * time.Second,
		Args: map[string]interface{}{"method":"MustOnceAfterWithFunArgs"},
		Action: controller.Index,
	})

	time.AfterFunc(3 * time.Second, func() {
		Timer.Clear(id3)
	})
	Timer.Clear(id3)
}
