// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"strings"

	"github.com/g41797/sputnik"
)

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

var states = [...]string{"Unknown", "Created", "Submitted", "Received", "InProcess", "Cancelled", "Finished", "Failed"}

func (jst JobState) String() string {
	return states[jst]
}

func (jst JobState) EnumIndex() int {
	return int(jst)
}

func BuildJobState(s string) JobState {
	s = strings.ToLower(s)
	for i, str := range states {
		if str == s {
			return JobState(i)
		}
	}

	return Unknown
}

type JobStatus interface {
	UID() string

	State() JobState

	// Any  information used by internal implementation
	Memento() any

	StatusTo(m sputnik.Msg) error
	StatusFrom(m sputnik.Msg) error
}

type JobAttribute struct {
	Name  string
	Value string
}

// Information for job submission
type JobOrder interface {
	// Name of job handler  responsible for Job processing
	Handler() string

	// List of attributes - may be empty
	Attributes() []JobAttribute

	// JSON byte array - may be nil/empty
	Payload() []byte

	OrderTo(m sputnik.Msg) error
	OrderFrom(m sputnik.Msg) error
}

// Job - unit of processing
type Job interface {
	UID() string
	JobOrder
	To(m sputnik.Msg) error
	From(m sputnik.Msg) error
}
