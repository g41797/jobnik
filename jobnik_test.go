// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik_test

import (
	"context"
	"fmt"
	"log"

	"github.com/g41797/jobnik"
)

// Factory function
func jobHandlerFactory() (jobnik.Jobnik, error) {
	return new(jobHandler), nil
}

// jobHandler supports 2 "commands" for jobnik: 'printattributes' and 'looptillcancel'
// so factory should be registered twice:
func init() {
	jobnik.RegisterFactory("printAttributes", jobHandlerFactory)
}

func init() {
	jobnik.RegisterFactory("loopTillCancel", jobHandlerFactory)
}

type subHandler func(ctx context.Context, job jobnik.Job) (jobnik.JobStatus, error)

type jobHandler struct {
	cnf  config
	hnds map[string]subHandler
}

func (jh *jobHandler) InitOnce(jsc string) error {

	jh.hnds = make(map[string]subHandler)
	jh.hnds["printattributes"] = jh.printAttributes
	jh.hnds["looptillcancel"] = jh.loopTillCancel

	var tempcnf config
	tempcnf.setDefault()

	if err := tempcnf.unmarshall(jsc); err != nil {
		jh.logError(err)
		return err
	}

	jh.cnf = tempcnf
	return nil
}

func (jh *jobHandler) FinishOnce() error {
	if jh.cnf.isLogInfo() {
		log.Default().Printf("jobHandler finished")
	}
	return nil
}

func (jh *jobHandler) Process(ctx context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	if job == nil {
		err := fmt.Errorf("empty job")
		jh.logError(err)
		return jobnik.JobStatus{}, err
	}

	actnm := job.Name()

	acthnd, exists := jh.hnds[actnm]

	if !exists {
		return jobnik.JobStatus{
			UID:      job.UID(),
			State:    jobnik.Failed,
			Addendum: fmt.Sprintf("%s is not supported", actnm)}, nil
	}

	return acthnd(ctx, job)
}

func (jh *jobHandler) printAttributes(ctx context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	return jobnik.JobStatus{}, fmt.Errorf("printAttributes is not implemented")
}

func (jh *jobHandler) loopTillCancel(ctx context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	return jobnik.JobStatus{}, fmt.Errorf("loopTillCancel is not implemented")
}

func (jh *jobHandler) logError(err error) {
	if err == nil {
		return
	}

	if !jh.cnf.isLogErrors() {
		return
	}

	log.Default().Print(err)
}
