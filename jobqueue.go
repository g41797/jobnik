// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/g41797/sputnik"
)

type Submitter interface {
	// Submit job to queue
	// Returns JobStatus for further job tracking
	Submit(jo JobOrder) (JobStatus, error)

	// Returns current job status
	// Lack of the job returned as Unknown JobState
	// error returned for failed access to job queue
	Check(js JobStatus) (JobStatus, error)
}

type Receiver interface {
	// Register receive callback
	OnReceive(rcv func(j Job)) error

	// Allows to receive next job from the queue
	// Not accumulated - 100 Acks allow 1 receive :-(
	// js may be nil (at lest it's nil for the first receive)
	// Informs queue about changed job status (for not nil)
	Ack(js JobStatus) error
}

type JobQueue interface {
	sputnik.ServerConnector
	Submitter
	Receiver
}

type JobQueueFactory func() (JobQueue, error)

// Store factory for further usage
// name of queue is stored in lower case
func RegisterJobQueueFactory(name string, fact JobQueueFactory) {
	if len(name) == 0 {
		log.Panic("empty JobQueue name")
	}
	if fact == nil {
		log.Panicf("nil JobQueue factory for %s", name)
	}

	if _, exists := factories.LoadOrStore(strings.ToLower(name), fact); exists {
		log.Panicf("JobQueue factory for %s already exists", name)
	}
	return
}

func CreateJobQueue(name string) (JobQueue, error) {

	fact, exists := factories.Load(strings.ToLower(name))
	if !exists {
		return nil, fmt.Errorf("factory for %s does not exist", name)
	}

	tr, err := fact.(JobQueueFactory)()

	return tr, err
}

var factories sync.Map
