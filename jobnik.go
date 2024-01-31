// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import "context"

// Flow of job processing:
//  Submitter: Creates Job and submits it to Queue
//	Worker: Listens to Queue, receives Job, activates corresponding Jobnik
//	Jobnik: Processes Job, returns JobStatus to Worker
//	Worker: Updates Queue with JobStatus
//	Anyone: Listens to Queue, receives JobStatus

// Jobnik is small part (plugin) of whole flow. It responsible for JobProcessing, that's all.
// But placement of Jobnik to own package  allows convenient development.
// BTW - You can use it on it's own
type Jobnik interface {
	// Initiation of Jobnik - once for lifecycle.
	// name - identification of Jobnik, allows to use the same code
	// for processing of different types of jobs.
	// jsc - JSON string with configuration.
	InitOnce(name string = nil, jsc string = nil) error

	// Stop processing, clean resources - once for lifecycle.
	FinishOnce() error

	// Process Job
	// cncl used for external cancel of running job.
	// how to use it see https://www.sohamkamani.com/golang/context/
	// error is returned only for wrong Job information, e.g. failure of
	// de-marshalling of JSON job payload
	// Failure or Cancel are valid states. It should be reflected in JobStatus.
	// error for these cases should be nil
	Process(cncl context.Context, Job job) (JobStatus, error)
}