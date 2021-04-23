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
	timerMap map[int64]timerStruct
}

type timerStruct struct {
	d time.Duration
	t *time.Timer
	action func(args ...map[string]interface{}) map[string]interface{}
}


var Timer timer
func init() {
	Timer = timer{
		id:       0,
		mutex:    sync.Mutex{},
		stopChan: make(chan int64),
		timerMap: make(map[int64]timerStruct),
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

//必须执行一次的定时器，无法销毁
func (t *timer) MustOnceAfter (second time.Duration,action func(args ...map[string]interface{}) map[string]interface{} ) {
	time.AfterFunc(second, func() {
		res := action()
		log.Printf("MustOnceAfter:%+v",res)
	})
}

//一次性定时器,可根据timerId提前销毁
func (t *timer) After (second time.Duration,action func(args ...map[string]interface{}) map[string]interface{}) int64 {
	timerId := t.newTimerId()

	t.timerMap[timerId] = timerStruct{
		d:      second,
		t:      time.AfterFunc(second, func() {
			res := action()
			log.Printf("After:%+v",res)
		}),
		action: nil,
	}
	return timerId
}



//周期定时器
func (t *timer) Loop (second time.Duration,action func(args ...map[string]interface{}) map[string]interface{},stop ...int) int64 {
	timerId := t.newTimerId()
	t.timerMap[timerId] = timerStruct{
		d:      second,
		t:      time.NewTimer(second),
		action: nil,
	}
	go func(id int64,s timerStruct) {
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
				t.timerMap[id].t.Reset(s.d)
				//s.Action(s.Args)
				res := action()
				log.Printf("Loop:%+v",res)
			}
		}
	}(timerId,t.timerMap[timerId])

	return timerId
}


func Call (args ...map[string]interface{}) map[string]interface{} {
	log.Printf("-------------- Call args:%+v --------------------",args)
	r := map[string]interface{}{"A":"B"}
	if len(args) > 0 {
		r = args[0]
	}
	return r
}

func test() {
	id := Timer.After(time.Second * 7 , Call)
	log.Printf("Timer After id:%+v",id)

	id2 := Timer.After(time.Second * 3 , Call)
	log.Printf("Timer Tick id:%+v",id2)

	id3 := Timer.Loop (time.Second * 1 , func(args ...map[string]interface{}) map[string]interface{} {
		return Call(map[string]interface{}{"method":"loop-call","test1":1})
	},1,5)
	log.Printf("Timer Tick id:%+v",id3)

	Timer.MustOnceAfter(time.Second * 5 , func(args ...map[string]interface{}) map[string]interface{} {
		return Call(controller.Index(map[string]interface{}{"method":"MustOnceAfter"}))
	})

	time.AfterFunc(3 * time.Second, func() {
		Timer.Clear(id3)
	})
	Timer.Clear(id)
}
