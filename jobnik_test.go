// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik_test

import (
	"context"
	"encoding/json"
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

//-------------------------------------------------------------------------------------------------------
// Useful JSON&Go links:
//
// https://jsonlint.com
// https://betterstack.com/community/guides/scaling-go/json-in-go/
// https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go#using-a-struct-to-generate-json
// https://blog.boot.dev/golang/json-golang/
// https://biscuit.ninja/posts/unmarshalling-json-with-null-booleans-in-go/#solution-1-use-pointers
//-------------------------------------------------------------------------------------------------------

// bool needs special attention
type config struct {
	IsLogErrors   *bool `json:"errors,omitempty"`
	IsLogWarnings *bool `json:"warnings,omitempty"`
	IsLogInfo     *bool `json:"info,omitempty"`
}

func (cnf *config) setDefault() {
	cnf.IsLogErrors = new(bool)
	*cnf.IsLogErrors = true
}

func (cnf *config) unmarshall(jstr string) error {

	if len(jstr) == 0 {
		return fmt.Errorf("empty JSON string")
	}

	cnf.setDefault()

	return json.Unmarshal([]byte(jstr), cnf)
}

func (cnf *config) isLogErrors() bool {
	if cnf.IsLogErrors == nil {
		return false
	}

	return *cnf.IsLogErrors
}

func (cnf *config) isLogWarnings() bool {
	if cnf.IsLogWarnings == nil {
		return false
	}

	return *cnf.IsLogWarnings
}

func (cnf *config) isLogInfo() bool {
	if cnf.IsLogInfo == nil {
		return false
	}

	return *cnf.IsLogInfo
}

type subHandler func(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error)

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

func (jh *jobHandler) Process(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
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

	return acthnd(cncl, job)
}

func (jh *jobHandler) printAttributes(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	return jobnik.JobStatus{}, fmt.Errorf("printAttributes is not implemented")
}

func (jh *jobHandler) loopTillCancel(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
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
