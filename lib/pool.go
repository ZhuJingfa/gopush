// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"sync"
	"log"

	apns "github.com/sideshow/apns2"

	loglocal "zooinit/log"
)

type Pool struct {
	//worker pool
	Workers       []Worker
	WorkerIDIndex int

	//worker poll size
	//MiniSpare <= now <= Capacity
	Size          int

	//pool capacity
	Capacity      int

	//mini spare worker
	MiniSpare     int

	//max spare worker
	MaxSpare      int

	//pool lock
	Lock          sync.Mutex

	wg            sync.WaitGroup

	OKLogger      *loglocal.BufferedFileLogger
	FailLogger    *loglocal.BufferedFileLogger

	Env           EnvInfo
}

// create a new worker pool
func NewPool(Size, Capacity, MiniSpare, MaxSpare int, Env EnvInfo) (*Pool, error) {
	workers := make([]Worker, Size, Capacity)
	WorkerIDIndex := 0

	pool := &Pool{Size:Size, Capacity:Capacity, MaxSpare:MaxSpare, MiniSpare:MiniSpare}

	for iter, _ := range workers {
		worker, err := Env.CreateWorker()
		if err != nil {
			return nil, err
		}

		//不能用append, 会增长数组.
		//workers=append(workers, worker)
		//start from 0
		worker.SetWorkerID(iter)
		WorkerIDIndex = iter

		worker.SetPool(pool)

		workers[iter] = worker
	}

	//WorkerIDIndex is last new one
	pool.WorkerIDIndex = WorkerIDIndex + 1
	pool.Workers = workers

	return pool, nil
}

func (p *Pool) FetchASpareWork() Worker {
	return nil
}

func (p *Pool) Run(list *DeviceQueue, msg *apns.Notification) {
	con, err := msg.MarshalJSON()
	if err != nil {
		p.Env.GetLogger().Println("msg.MarshalJSON() found error:", err)
		return
	}
	p.Env.GetLogger().Println("Receive new push task: " + string(con))

	if len(p.Workers) != p.Size {
		p.Env.GetLogger().Fatalln("Found exception of pool: len(p.Workers)!=p.Size: ", len(p.Workers), p.Size)
	}

	// start up worker
	for _, worker := range p.Workers {
		p.wg.Add(1)
		//env.GetLogger().Println(worker.GetWorkerName()+" ...")

		//TODO Here has an error Mode if run in anonymous func, worker started is not in expected mode
		go func() {
			worker.Run()
			p.wg.Done()
		}()

		//goroutine wg.done in worker.Run()
	}

	//p.Test(list, msg)
	p.Send(list, msg)

	// wait for all done worker.Run() / worker.Subscribe()
	p.wg.Wait()
}

func (p *Pool) Send(list *DeviceQueue, msg *apns.Notification) {
	// Queue data publish
	go list.Publish()

	for _, worker := range p.Workers {

		p.wg.Add(1)
		go worker.Subscribe(list, msg)

		//goroutine wg.done in worker.Subscribe()
	}
}

func (p *Pool) GetOKLogger() (*loglocal.BufferedFileLogger) {
	if p.OKLogger == nil {
		p.OKLogger = p.getInternalLogger("ok")
	}

	return p.OKLogger
}

func (p *Pool) GetFailLogger() (*loglocal.BufferedFileLogger) {
	if p.FailLogger == nil {
		p.FailLogger = p.getInternalLogger("fail")
	}

	return p.FailLogger
}

func (p *Pool) getInternalLogger(logtype string) (*loglocal.BufferedFileLogger) {
	filename := loglocal.GenerateFileLogPathName(p.Env.GetLogPath(), "pool_" + logtype)
	file, err := loglocal.NewFileLog(filename)
	if err != nil {
		log.Fatalln(err)
	}

	logger := log.New(file, "", log.Ldate | log.Ltime | log.Lmicroseconds) // add time for stat
	return loglocal.GetBufferedFileLogger(file, logger)
}