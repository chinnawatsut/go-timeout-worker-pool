package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type ProcessorFunc func(resource interface{})

type Job struct {
	id       int
	Resource interface{}
}

type Pool struct {
	numRoutines   int
	numJobs       int
	jobs          chan Job
	resultCounter chan bool
	done          chan bool
}

func NewPool(numRoutines int) *Pool {
	r := &Pool{numRoutines: numRoutines}
	r.jobs = make(chan Job, numRoutines)
	r.resultCounter = make(chan bool, numRoutines)
	return r
}

func (m *Pool) Start(resources []interface{}, timeOutSec int64, procFunc ProcessorFunc) {
	m.numJobs = len(resources)
	m.done = make(chan bool)
	go m.allocate(resources)
	go m.counter()
	go m.workerPool(procFunc)
	timeOut := time.After(time.Duration(timeOutSec) * time.Second)

	select {
	case <-m.done:
		fmt.Println("COMPLETE")
		return
	case <-timeOut:
		fmt.Println("COMPLETE WITH TIMEOUT")
		return
	}
}

func (m *Pool) allocate(jobs []interface{}) {
	defer close(m.jobs)
	for i, v := range jobs {
		job := Job{id: i, Resource: v}
		m.jobs <- job
		fmt.Println("ADD JOB:", i)
	}
}

func (m *Pool) work(wg *sync.WaitGroup, processor ProcessorFunc, workerNumber int) {
	defer func() {
		fmt.Println("Worker No.", workerNumber, "is stop")
	}()
	defer wg.Done()
	fmt.Println("Worker No.", workerNumber, "is working")
	for job := range m.jobs {
		fmt.Println(workerNumber, ": RUN")
		processor(job.Resource)
		m.resultCounter <- true
		fmt.Println(workerNumber, ": DONE")
	}

}

func (m *Pool) workerPool(processor ProcessorFunc) {
	defer close(m.resultCounter)
	var wg sync.WaitGroup
	for i := 0; i < m.numRoutines; i++ {
		wg.Add(1)
		fmt.Println("ADD WOKER No.", i)
		go m.work(&wg, processor, i)
	}
	fmt.Println("-----")
	wg.Wait()
}

func (m *Pool) counter() {
	for x := 0; x < m.numJobs; x++ {
		<-m.resultCounter
	}
	m.done <- true
}
