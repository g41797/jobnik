package jobnik_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/g41797/jobnik"
)

type longJobHandler struct {
	loops  int `json:"LoopsCount"`
	slpmls int `json:"SleepMlsec"`
}

// Factory function
func longJobHandlerFactory() (jobnik.Jobnik, error) {
	return new(longJobHandler), nil
}

// Register handler:
func init() {
	jobnik.RegisterFactory("ProcessLongJob", longJobHandlerFactory)
}

// Works without initial configuration
func (jh *longJobHandler) InitOnce(_ string) error {
	jh.setDefaults()
	return nil
}

func (jh *longJobHandler) setDefaults() {
	jh.loops = 100
	jh.slpmls = 1
	return
}

func (jh *longJobHandler) FinishOnce() error {
	return nil
}

func (jh *longJobHandler) Process(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	if job == nil {
		err := fmt.Errorf("empty job")
		return jobnik.JobStatus{}, err
	}

	payload := job.Payload()
	if len(payload) == 0 {
		err := fmt.Errorf("empty job payload")
		return jobnik.JobStatus{}, err
	}

	if err := json.Unmarshal([]byte(payload), jh); err != nil {
		jh.setDefaults()
		return jobnik.JobStatus{}, err
	}
	return jobnik.JobStatus{}, fmt.Errorf("ProcessLongJob is not implemented yet")
}
