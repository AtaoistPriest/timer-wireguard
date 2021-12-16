package main

import (
	"log"
	"sync"
	"time"
)

type SA struct {

	KeepAlive 		Timer
	KeyExchange 	Timer
	WRLock    		sync.RWMutex
	sessions    	map[[64]byte]sess
}

type sess struct {
	TaskTimer 	Timer
	WRLock    	sync.RWMutex
}

func (sa *SA)Disconnect(){
	log.Print("over")
}

const(
	RETRANSMIT = 0
	ALWAYSDONE = 1
)

type Timer struct {
	Clock		*time.Timer
	RWLock		sync.RWMutex
	isAlive		bool
	Task		TimerTask
}

type TimerTask struct {
	Type 		int
	Data 		string
	Id 			string
	TaskCnt 	int
	MaxCnt		int
	Duration 	time.Duration
}

func (timer *Timer) ReTry(){
	timer.RWLock.Lock()
	defer timer.RWLock.Unlock()
	timer.isAlive = true
	if timer.Task.Type == RETRANSMIT{
		timer.Task.TaskCnt++
	}
	timer.Clock.Reset(timer.Task.Duration)
}

func (timer *Timer) Reset(){
	timer.RWLock.Lock()
	defer timer.RWLock.Unlock()
	timer.isAlive = true
	timer.Task.TaskCnt = 0
	timer.Clock.Reset(timer.Task.Duration)
}

func (timer *Timer) Close(){
	timer.RWLock.Lock()
	defer timer.RWLock.Unlock()
	timer.isAlive = false
	timer.Task.TaskCnt = 0
	timer.Clock.Stop()
	return
}

func (timer *Timer) IsAlive() bool{
	timer.RWLock.Lock()
	defer timer.RWLock.Unlock()
	if timer.Task.Type == RETRANSMIT && timer.Task.TaskCnt >= timer.Task.MaxCnt{
		timer.isAlive = false
	}
	return timer.isAlive
}

func (timer *Timer)Start(handleRetry func(task TimerTask), handleOver func(task TimerTask)){
	timer.isAlive = true
	timer.Clock = time.AfterFunc(timer.Task.Duration, func() {
		if timer.IsAlive(){
			if handleRetry != nil{
				handleRetry(timer.Task)
			}
			timer.ReTry()
		}else{
			if handleOver != nil{
				handleOver(timer.Task)
			}
			timer.Close()
		}
	})
}

func taskExecute(task TimerTask){
	log.Print("do something. It's saying : ", task.Data)
}

func taskOver(task TimerTask){
	log.Print("Over. It's saying : ", task.Data)
}

func main(){

	log.Print("begin")
	tmp := Timer{
		Task: TimerTask{
			Type: RETRANSMIT,
			Data: "I'm here!ha ha",
			Duration: time.Second * 3,
			Id: "IDIDIDIDDIDID",
			TaskCnt: 0,
			MaxCnt: 3,
		},
	}
	sa := SA{
		KeepAlive: tmp,
	}
	sa.KeepAlive.Start(taskExecute, nil)
	time.Sleep(time.Hour)
}
