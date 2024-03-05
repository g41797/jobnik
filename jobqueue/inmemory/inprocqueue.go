// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/g41797/jobnik"
	"github.com/g41797/kissngoqueue"
	"github.com/g41797/sputnik"
	"github.com/google/uuid"
)

const INMEMORYQUEUE = "inmemoryqueue"

func init() {
	jobnik.RegisterJobQueueFactory(INMEMORYQUEUE, queueFactory)
}

func queueFactory() (jobnik.JobQueue, error) {
	q := new(inprocqueue)
	q.SetState(true)
	q.q = kissngoqueue.NewQueue[jobnik.Job]()
	return q, nil
}

var _ jobnik.JobQueue = &inprocqueue{}

type inprocqueue struct {
	lock sync.Mutex
	sputnik.DummyConnector
	q      *kissngoqueue.Queue[jobnik.Job]
	states sync.Map
}

func (q *inprocqueue) Submit(jo jobnik.JobOrder) (jobnik.JobStatus, error) {

	if !q.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	job := new(jobnik.DefaultJob)

	job.Name = jo.Handler()
	copy(job.Atrbs, jo.Attributes())
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

func (q *inprocqueue) Receive(ctx context.Context) (jobnik.Job, error) {
	if !q.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	return nil, nil
}

func (q *inprocqueue) Ack(js jobnik.JobStatus) error {

	if !q.IsConnected() {
		return fmt.Errorf("not connected")
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
