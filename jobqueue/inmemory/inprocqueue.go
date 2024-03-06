// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package inmemory

import (
	"fmt"
	"sync"

	"github.com/g41797/jobnik"
	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
	"github.com/google/uuid"
)

func init() {
	jobnik.RegisterJobQueueFactory(jobnik.INMEMORYQUEUE, queueFactory)
}

func queueFactory() (jobnik.JobQueue, error) {
	q := new(inprocqueue)
	q.SetState(true)
	q.q = kissngoqueue.NewQueue[jobnik.Job]()
	return q, nil
}

var _ jobnik.JobQueue = &inprocqueue{}

// Non production code, used for the development/testing within one process
type inprocqueue struct {
	lock sync.Mutex
	sputnik.DummyConnector
	q      *kissngoqueue.Queue[jobnik.Job]
	states sync.Map
	rcv    func(j jobnik.Job)
	trg    chan struct{}
}

func (q *inprocqueue) Submit(jo jobnik.JobOrder) (jobnik.JobStatus, error) {

	if !q.IsConnected() {
		q.stopRecv()
		return nil, fmt.Errorf("not connected")
	}

	job := new(jobnik.DefaultJob)

	job.Name = jo.Handler()

	job.Atrbs = make([]jobnik.JobAttribute, len(jo.Attributes()))
	copy(job.Atrbs, jo.Attributes())

	job.Pld = make([]byte, len(jo.Payload()))
	copy(job.Pld, jo.Payload())

	uid := uuid.New().String()
	job.Uid = uid
	job.St = jobnik.Unknown

	q.lock.Lock()
	defer q.lock.Unlock()

	if ok := q.q.Put(job); !ok {
		return nil, fmt.Errorf("putmt to queue failed")
	}

	js := new(jobnik.DefaultJobStatus)
	js.Uid = uid
	js.St = jobnik.Submitted
	js.Mem = nil

	q.states.Store(js.Uid, js.St)

	return js, nil
}

func (q *inprocqueue) Check(js jobnik.JobStatus) (jobnik.JobStatus, error) {

	if !q.IsConnected() {
		q.stopRecv()
		return nil, fmt.Errorf("not connected")
	}

	jst := new(jobnik.DefaultJobStatus)
	jst.Uid = js.UID()
	jst.St = jobnik.Unknown
	jst.Mem = nil

	state, exists := q.states.Load(js.UID())

	if exists {
		jst.St = state.(jobnik.JobState)
	}

	return jst, nil
}

func (q *inprocqueue) OnReceive(rcv func(j jobnik.Job)) error {

	if !q.IsConnected() {
		q.stopRecv()
		return fmt.Errorf("not connected")
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	q.rcv = rcv

	return nil
}

func (q *inprocqueue) Ack(js jobnik.JobStatus) error {

	if !q.IsConnected() {
		q.stopRecv()
		return fmt.Errorf("not connected")
	}

	if err := q.allowRecv(); err != nil {
		return err
	}

	if js == nil {
		return nil
	}

	state, exists := q.states.Load(js.UID())

	if !exists {
		return nil
	}

	if state == jobnik.Finished {
		q.states.Delete(js.UID())
		return nil
	}

	if state != js.State() {
		q.states.Store(js.UID(), js.State())
	}

	return nil
}

func (q *inprocqueue) allowRecv() error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.rcv == nil {
		return fmt.Errorf("wrong flow")
	}

	if q.trg != nil {
		return nil
	}

	q.trg = make(chan struct{}, 1)

	go q.waitRecv()

	return nil
}

func (q *inprocqueue) stopRecv() {
	if q.trg != nil {
		close(q.trg)
		q.trg = nil
	}
}

func (q *inprocqueue) waitRecv() {
	job, _ := q.q.WaitMT(q.trg)
	q.lock.Lock()
	defer q.lock.Unlock()
	q.stopRecv()

	if (q.rcv == nil) || job == nil {
		return
	}

	q.states.Store(job.UID(), jobnik.InProcess)

	q.rcv(job)

	return
}

func (q *inprocqueue) Stop() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.stopRecv()

	q.Disconnect()

	return
}
