// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

type JobState int

const (
	Unknown JobState = iota

	// All information was added, ready for submit to queue
	Created JobState = iota + 1

	// Within Queue
	Submitted JobState = iota + 2

	// Received on Worker side
	Received JobState = iota + 3

	// Processed by jobnik
	InProcess JobState = iota + 4

	// Cancelled by user or another application
	Cancelled JobState = iota + 5

	// Well done
	Finished JobState = iota + 6

	// Obvious
	Failed JobState = iota + 7
)

type JobStatus struct {
	UID   string
	State JobState

	// Any  useful information,
	// e.g. initiator of Cancel, or reason of failure
	Addendum string
}

type JobAttribute struct {
	Name  string
	Value string
}

// Pay attention - you cannot get JobStatus directly from
// Job itself
type Job interface {
	// Unique Job ID
	UID() string

	// Name of jobnik(handler) responsible for Job processing
	Name() string

	// List of attributes - may be empty
	Attributes() []JobAttribute

	// JSON string - may be empty
	Payload() string
}

func NewJob(id string, name string, atrbs []JobAttribute, payload string) Job {
	return &defaultJob{id, name, atrbs, payload}
}

type defaultJob struct {
	id      string
	name    string
	atrbs   []JobAttribute
	payload string
}

func (job *defaultJob) UID() string { return job.id }

func (job *defaultJob) Name() string { return job.name }

func (job *defaultJob) Attributes() []JobAttribute { return job.atrbs }

func (job *defaultJob) Payload() string { return job.payload }

func (jst JobState) String() string {
	return []string{"Unknown", "Created", "Submitted", "Received", "InProcess", "Cancelled", "Finished", "Failed"}[jst]
}