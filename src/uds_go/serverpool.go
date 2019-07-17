package main

import (
	"fmt"
	"strconv"
	"time"
)

type Job struct {
	Payload *Payload
	requestType string
	jobName string
}

type Payload struct {
	dbRequest DbRequest
	key string
	table string
}

type Worker struct {
	WorkerName string
	WorkerPool chan WorkerWrapper
	WorkerWrapper WorkerWrapper
	Client Client
}

type WorkerWrapper struct {
	JobChannel chan *Job
	RetChannel chan interface{}
	ItrChannel chan bool
	Quit chan bool
}

func NewWorker(workerPool chan WorkerWrapper, name string, client Client) *Worker {
	fmt.Println("create worker:", name)
	return &Worker{WorkerName:name,
		WorkerPool: workerPool,
		WorkerWrapper: WorkerWrapper{JobChannel:make(chan *Job), RetChannel:make(chan interface{}), ItrChannel:make(chan bool, 1), Quit:make(chan bool),},
		Client: client,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.WorkerWrapper
			fmt.Println(w.WorkerName, "register to pool")
			select {
			case job := <- w.WorkerWrapper.JobChannel:
				fmt.Println(w.WorkerName, "recieved a job, current work pool size is", len(w.WorkerPool))
				var result interface{}
				if job.requestType == "get" {
					rsp, err := getRequest(job.Payload.dbRequest, w.Client, job.Payload.key, job.Payload.table)
					if err != nil {
						result = fmt.Sprintf("[%s] to get [%s] failed", w.WorkerName, job.Payload.key)
					} else {
						fmt.Println(w.WorkerName, "get result", rsp)
						result = rsp
					}
				} else if job.requestType == "put" {
					err := putRequest(job.Payload.dbRequest, w.Client, job.Payload.table)
					if err != nil {
						result = fmt.Sprintf("[%s] put failed", w.WorkerName)
					} else {
						result = fmt.Sprintf("[%s] put successed", w.WorkerName)
					}
				} else {
					result = fmt.Sprintf("[%s] request unsupported", w.WorkerName)
				}
					select {
					case w.WorkerWrapper.RetChannel <- result:
					case <- w.WorkerWrapper.ItrChannel:
						fmt.Println(w.WorkerName, "interrupted")
					}
			case <- w.WorkerWrapper.Quit:
				fmt.Println("close client", w.WorkerName)
				w.Client.close()
				return
			}
		}
	}()
}

func (w *Worker) Stop(){
	go func(){
		w.WorkerWrapper.Quit <- true
	}()
}

type Dispatcher struct {
	dispatcherType string
	endpoints string
	maxWorkers int
	WorkerPool chan WorkerWrapper
}

func NewDispatcher(maxWorkers int, dType string, addresses string) *Dispatcher {
	pool := make(chan WorkerWrapper, maxWorkers)
	return &Dispatcher{
		WorkerPool:pool,
		dispatcherType: dType,
		maxWorkers:maxWorkers,
		endpoints:addresses,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		if d.dispatcherType == "riak" {
			client := &RiakClient{}
			client.setEndpoint(d.endpoints)
			if err := client.build(); err != nil {
				panic(fmt.Sprintf("cannot build riak client: %v", err))
			}
			if err := client.open(); err != nil {
				panic(fmt.Sprintf("cannot open riak client: %v", err))
			}
			worker := NewWorker(d.WorkerPool, fmt.Sprintf("riak-worker-%s",strconv.Itoa(i)), client)
			worker.Start()
		} else if d.dispatcherType == "hbase" {
			client := &HBaseClient{}
			client.setEndpoint(d.endpoints)
			if err := client.build(); err != nil {
				panic(fmt.Sprintf("cannot build hbase client: %v", err))
			}
			worker := NewWorker(d.WorkerPool, fmt.Sprintf("hbase-work-%s",strconv.Itoa(i)), client)
			worker.Start()
		} else {
			panic(fmt.Sprintf("not support dispatcher type: %s", d.dispatcherType))
		}
	}
	for {
		worker_len := len(d.WorkerPool)
		if worker_len != d.maxWorkers {
			fmt.Println("wait for 1 sec to check if all workers started", worker_len)
			time.Sleep(1000 * time.Millisecond)
		} else {
			fmt.Println("All Worker Started", worker_len)
			break
		}
	}
}


func (d *Dispatcher) Execute(job *Job, timeout int) interface{} {
	wrapper := <- d.WorkerPool
	wrapper.JobChannel <- job

	tout := time.NewTimer(time.Duration(timeout) * time.Millisecond)

	var rep interface{}
	select {
	case res := <- wrapper.RetChannel:
		rep = res
	case <- tout.C:
		wrapper.ItrChannel <- true
		rep = "{}"
	}
	tout.Stop()
	return rep
}

func (d *Dispatcher) Exit() {
	for i := 0; i < d.maxWorkers; i++ {
		wrapper := <- d.WorkerPool
		wrapper.Quit <- true
	}
}

var dRiak *Dispatcher
var dHBase *Dispatcher

func initialize() {
	maxWorkers := 10
	dRiak = NewDispatcher(maxWorkers, "riak", "qa-database-riak-userdata-xv-02.xv.dc.openx.org:8087")
	dRiak.Run()
	dHBase = NewDispatcher(maxWorkers, "hbase", "hbase-userprof-xv-01.xv.dc.openx.org:2181,hbase-userprof-xv-03.xv.dc.openx.org:2181,hbase-userprof-xv-06.xv.dc.openx.org:2181")
	dHBase.Run()
}

func getRiakDispatcher() *Dispatcher {
	return dRiak
}

func getHBaseDispatcher() *Dispatcher {
	return dHBase
}



