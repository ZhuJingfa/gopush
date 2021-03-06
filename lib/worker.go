// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import "sync"

const (
	WORKER_STATUS_SPARE = iota
	WORKER_STATUS_RUNNING

	WORKER_COMMAND_SEND = iota
	WORKER_COMMAND_STOP
)

type Worker interface {
	Run()

	Stop()

	// Subscribe device goroutine run
	Subscribe(task *Task)

	// Push a message
	Push(msg MessageInterface, Device string) (*WorkerResponse)

	// Fetch a worker identify
	GetWorkerName() (string)

	// Set worker id
	SetWorkerID(id int) (bool)
	GetWorkerID() (int)

	// Set worker's belonging pool
	SetPool(pool *Pool) (bool)
	GetPool() (*Pool)

	Destroy() (error)

	GetStatus() int

	//worker's uuid identify
	GetUUID() string

	GetLockPtr() *sync.Mutex
}

type WorkerRequeset struct {
	Message MessageInterface

	//the specified device send to
	Device  string

	Cmd     int
}

//Create a new request
func NewWorkerRequeset(msg MessageInterface, device string, cmd int) *WorkerRequeset {
	return &WorkerRequeset{Message:msg, Device:device, Cmd:cmd}
}

type WorkerResponse struct {
	Response interface{}

	//the specified device send to
	Device   string

	Error    error
}
