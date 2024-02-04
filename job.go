// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"encoding/json"
	"fmt"
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

func (jst JobState) String() string {
	return []string{"Unknown", "Created", "Submitted", "Received", "InProcess", "Cancelled", "Finished", "Failed"}[jst]
}

type JobStatus struct {
	UID   string
	State JobState

	// Any  useful information,
	// e.g. initiator of Cancel, or reason of failure
	Addendum string
}

func (jst *JobStatus) String() string {
	if len(jst.Addendum) == 0 {
		return fmt.Sprintf(" ID %s State %s ", jst.UID, jst.State.String())
	}
	return fmt.Sprintf(" ID %s State %s Additional information %s", jst.UID, jst.State.String(), jst.Addendum)
}

type JobAttribute struct {
	Name  string
	Value string
}

// Information for job submission
type JobOrder interface {
	// Name of jobnik(handler) responsible for Job processing
	Name() string

	// List of attributes - may be empty
	Attributes() []JobAttribute

	// JSON string - may be empty
	Payload() string
}

// Job - unit of processing
// Pay attention - you cannot get JobStatus directly from
// Job itself
type Job interface {
	// Unique Job ID created during submission
	UID() string

	JobOrder
}

// NewJobOrder returns error for empty name and failed JSON unmarshall of payload
// You can disable JSON validation assign true to 'dontval'
func NewJobOrder(name string, atrbs []JobAttribute, payload string, dontval bool) (JobOrder, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("empty name")
	}

	if !dontval && len(payload) > 0 {
		var unmarsh map[string]any
		if err := json.Unmarshal([]byte(payload), &unmarsh); err != nil {
			return nil, err
		}
	}

	return &defaultJobOrder{name, atrbs, payload}, nil
}

// NewJob returns error for empty id and/or name
// payload (as a rule JSON string) - is not validated
func NewJob(id string, name string, atrbs []JobAttribute, payload string) (Job, error) {
	if len(id) == 0 || len(name) == 0 {
		return nil, fmt.Errorf("wrong id or name")
	}
	return &defaultJob{id, defaultJobOrder{name, atrbs, payload}}, nil
}

// NewJobForOrder returns error for empty id and jo
func NewJobForOrder(id string, jo JobOrder) (Job, error) {
	if jo == nil {
		return nil, fmt.Errorf("job order")
	}
	return NewJob(id, jo.Name(), jo.Attributes(), jo.Payload())
}
