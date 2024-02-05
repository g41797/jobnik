package jobnik_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

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
	return jobnik.JobStatus{}, fmt.Errorf("ProcessLongJob is not implemented yet")
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
}
