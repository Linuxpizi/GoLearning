package sche

import (
	"sync"
	"time"
)

type taskFunc func()
type schedule struct {
	lock     sync.RWMutex //ptotect stopTask
	stopTask bool         //是否开启任务
	duration uint8        //运行的周期 单位s
	timeout  uint8        //如果stop，但是没有start多久后自动start 单位s
	start    taskFunc     //start需要执行的回调函数
	stop     taskFunc     //清理工作
}

// Scheduler ..
type Scheduler interface {
	StartTask()
	StopTask()
}

// NewTask ..
func NewTask(duration, timeout uint8, start, stop taskFunc) Scheduler {
	s := &schedule{
		stopTask: false, //默认服务启动
		duration: duration,
		timeout:  timeout,
		start:    start,
		stop:     stop,
	}
	go s.run()
	return s
}

func (s *schedule) run() {
	for {
		select {
		case <-time.Tick(time.Second * time.Duration(s.duration)):

			s.lock.Lock()
			if !s.stopTask {
				s.start()
			}
			s.lock.Unlock()
		}
	}
}

func (s *schedule) StartTask() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stopTask = false
}

func (s *schedule) StopTask() {
	s.lock.Lock()
	s.stopTask = true
	s.lock.Unlock()
	go func() {
		for {
			select {
			case <-time.After(time.Second * time.Duration(s.timeout)):
				s.lock.Lock()
				if s.stopTask {
					s.stopTask = false
					s.lock.Unlock()
					return
				}
				s.lock.Unlock()
			}
		}
	}()
}
