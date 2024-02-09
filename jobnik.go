// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"context"
	"fmt"
)

// Flow of job processing:
//	Submitter: Creates Job and submits it to Queue
//	Worker: Listens to Queue, receives Job, activates corresponding Jobnik
//	Jobnik: Processes Job, returns JobStatus to Worker
//	Worker: Updates Queue with JobStatus
//	Anyone: Listens to Queue, receives JobStatus

// Jobnik is small part (plugin) of whole flow. It's responsible for JobProcessing, that's all.
// But placement of Jobnik to own package  allows convenient development.
// BTW - You can use it on it's own
type Jobnik interface {
	// Initiation of Jobnik - once for lifecycle.
	// jsc (may be empty) - JSON string with configuration.
	InitOnce(jsc string) error

	// Stop processing, clean resources - once for lifecycle.
	FinishOnce() error

	// Process Job
	//
	// ctx used for external cancel of running job.
	// how to use it see:
	// https://www.sohamkamani.com/golang/context/
	// https://www.willem.dev/articles/context-cancellation-explained/
	// https://www.willem.dev/articles/context-todo-or-background/
	// If ctx is not initialized, context.Background() will be used
	//
	// error is returned only for wrong arguments, e.g. failure of
	// de-marshalling of JSON job payload, for this case content of returned JobStatus does not matter.
	// Failure or Cancel are valid states. It should be reflected in JobStatus.
	// error for these cases should be nil
	Process(ctx context.Context, job Job) (JobStatus, error)
}

// Don't create jobnik direct and don't export it from the package.
// Use NewJobnik function.
// name should be unique within process.
// name is non case sensitive: "CopyFiles" and "COPYFILES" are the same name
// (actually "copyfiles")
func NewJobnik(name string) (Jobnik, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("empty jobnik name")
	}

	jbn := new(guard)
	return jbn.tryCreate(name)
}

// In order to allow indirect creation:
// 1 - 	JobnikFactory should be provided for every jobnik in the process
// 2 - 	JobnikFactory should be registered before creation. Usually it will be
//		done within init()
//
//	 func init() {
//			jobnik.RegisterFactory("CopyFiles", cpfFactory)
//		}

// Provided by jobnik developer
type JobnikFactory func() (Jobnik, error)

// Store factory for further usage by NewJobnik
func RegisterFactory(name string, fact JobnikFactory) {
	storeFactory(name, fact)
}
