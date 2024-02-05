package jobnik_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/g41797/jobnik"
)

const longJobJobnikName = "ProcessLongJob"

type longJobHandler struct {
	Loops  int `json:"LoopsCount"`
	Slpmls int `json:"SleepMlsec"`
}

// Factory function
func longJobHandlerFactory() (jobnik.Jobnik, error) {
	return new(longJobHandler), nil
}

// Register handler:
func init() {
	jobnik.RegisterFactory(longJobJobnikName, longJobHandlerFactory)
}

// Works without initial configuration
func (jh *longJobHandler) InitOnce(_ string) error {
	jh.setDefaults()
	return nil
}

func (jh *longJobHandler) setDefaults() {
	jh.Loops = 100
	jh.Slpmls = 1
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

	for i := 0; i < jh.Loops; i++ {

		select {
		// see https://stackoverflow.com/questions/17573190/how-to-multiply-duration-by-integer
		case <-time.After(time.Duration(jh.Slpmls) * time.Millisecond):
			continue
		case <-cncl.Done():
			jst := jobnik.JobStatus{
				UID:      job.UID(),
				State:    jobnik.Cancelled,
				Addendum: fmt.Sprintf("Done %d loops", i)}
			return jst, nil
		}

	}

	return jobnik.JobStatus{job.UID(), jobnik.Finished, fmt.Sprintf("Done %d loops", jh.Loops)}, nil
}

//----------------------------------------------
// Useful links for using embedded folders/files
//----------------------------------------------
// https://pkg.go.dev/embed
// https://gobyexample.com/embed-directive
// https://www.iamyadav.com/blogs/a-guide-to-embedding-static-files-in-go

//go:embed testdata/wrong.json
var wj string

func TestWrongJSON(t *testing.T) {
	if _, err := jobnik.NewJobOrder("name", nil, wj, false); err == nil {
		t.Errorf("error == nil for wrong json")
	}
}

//go:embed testdata/longjobpayload1.json
var ljp1 string

func TestRightJSON(t *testing.T) {
	if _, err := jobnik.NewJobOrder("name", nil, ljp1, false); err != nil {
		t.Errorf("error == nil for right json")
	}
}

func TestProcess(t *testing.T) {
	if _, err := jobnik.NewJobnik("UNKNOWN"); err == nil {
		t.Errorf("error == nil for non existing jobnik")
	}

	job, _ := jobnik.NewJob("id", longJobJobnikName, nil, ljp1)

	jbnk, err := jobnik.NewJobnik(job.Name())
	if err != nil {
		t.Errorf("error != nil for existing jobnik %s", job.Name())
	}

	if err = jbnk.InitOnce(""); err != nil {
		t.Errorf("error %s for InitOnce", err.Error())
	}
	defer jbnk.FinishOnce()

	jbst, err := jbnk.Process(context.Background(), job)

	if err != nil {
		t.Errorf("error %s for Process", err.Error())
	}

	if jbst.State != jobnik.Finished {
		t.Errorf("wrong job state %s ", jbst.State)
	}

	//----------------------------------------------
	// Useful links for context cancel
	//----------------------------------------------
	// https://www.technicalfeeder.com/2023/01/golang-how-to-differentiate-between-context-cancel-and-timeout/
	// https://stackoverflow.com/questions/70042646/best-practices-on-go-context-cancelation-functions

	ctx, cancel := context.WithCancel(context.Background())
	recv := make(chan jobnik.JobStatus)
	procInLoop := func(done chan jobnik.JobStatus) {
		jbst, _ := jbnk.Process(ctx, job)
		done <- jbst
	}

	go procInLoop(recv)

	time.Sleep(time.Duration(1000) * time.Millisecond)

	cancel()

	jbst = <-recv

	if jbst.State != jobnik.Cancelled {
		t.Errorf("wrong job state %s ", jbst.State)
	}

}
