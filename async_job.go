package common

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
)

var AsyncJobIns *AsyncJob

func InitAysncJob(workerNum int) {
	AsyncJobIns = NewAsyncJob(workerNum)
}

// AsyncJob 异步写入, 考虑更通用的方式
type AsyncJob struct {
	ffs       chan func() error
	stop      chan struct{}
	workerNum int
	isRunning int32
}

func NewAsyncJob(workerNum int) *AsyncJob {
	af := &AsyncJob{
		ffs:       make(chan func() error, 1024),
		stop:      make(chan struct{}),
		workerNum: workerNum,
	}
	if workerNum <= 0 {
		af.workerNum = 1
	}
	af.Run()

	return af
}

// cleanBuf
func (af *AsyncJob) cleanBuf() {
	for {
		select {
		case ff := <-af.ffs:
			ff()
		default:
			return
		}
	}
}

func (af *AsyncJob) Run() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 0, 1) {
		return
	}
	for i := 0; i < af.workerNum; i++ {
		go af.runOuter()
	}
}

func (af *AsyncJob) runOuter() {
	for {
		select {
		case <-af.stop:
			af.stop <- struct{}{}
			return
		default:
			af.runrun()
		}
	}

}

func (af *AsyncJob) runrun() {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 20
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Printf("AsyncJob panic=%v\n%s\n", err, buf)
		}
	}()
	for {
		select {
		case ff, ok := <-af.ffs:
			if !ok {
				fmt.Fprintln(os.Stderr, "buf channel has been closed.")
				return
			}
			ff()
		case <-af.stop:
			af.cleanBuf()
			af.stop <- struct{}{}
			return
		}
	}

}

func (af *AsyncJob) AddJob(ff func() error) error {
	select {
	case af.ffs <- ff:
		return nil
	default:
		// warn write loss
		return fmt.Errorf("%s", "AsyncJob overflow")

	}
}

func (af *AsyncJob) AddJobNoLoss(ff func() error) {
	af.ffs <- ff
}

// Close safe clean
func (af *AsyncJob) Close() {
	if !atomic.CompareAndSwapInt32(&af.isRunning, 1, 0) {
		return
	}
	af.stop <- struct{}{}
	<-af.stop
}
