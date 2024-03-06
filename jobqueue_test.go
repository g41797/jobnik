// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik_test

import (
	"bytes"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/g41797/jobnik"
	_ "github.com/g41797/jobnik/jobqueue/inmemory"
	"github.com/g41797/kissngoqueue"
)

func newTester(name string) (*jqTest, error) {

	jobQueue, err := jobnik.CreateJobQueue(name)

	if err != nil {
		return nil, err
	}

	tst := new(jqTest)
	tst.jq = jobQueue
	tst.jobs = kissngoqueue.NewQueue[jobnik.Job]()
	return tst, nil
}

type jqTest struct {
	lock sync.Mutex
	jq   jobnik.JobQueue
	jobs *kissngoqueue.Queue[jobnik.Job]
}

func (tst *jqTest) onRecv(job jobnik.Job) {
	tst.jobs.PutMT(job)
}

func Test_FirstSteps(t *testing.T) {
	tst, err := newTester(jobnik.INMEMORYQUEUE)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(tst.jq.Stop)

	jord := &jobnik.DefaultJobOrder{
		Name:  "Any",
		Atrbs: []jobnik.JobAttribute{{"A1", "V1"}},
		Pld:   []byte("1234567890"),
	}

	if err = tst.jq.OnReceive(tst.onRecv); err != nil {
		t.Fatal(err)
	}

	jst, err := tst.jq.Submit(jord)

	if err != nil {
		t.Fatal(err)
	}

	if jst.State() != jobnik.Submitted {
		t.Fatalf("wrong job status %s", jst.State().String())
	}

	sjs, err := tst.jq.Check(jst)
	if err != nil {
		t.Fatal(err)
	}

	if jst.UID() != sjs.UID() {
		t.Fatalf("actual UID %s expected UID %s", sjs.UID(), jst.UID())
	}

	pjst := &jobnik.DefaultJobStatus{
		Uid: jst.UID(),
		St:  jobnik.Finished,
	}

	if err = tst.jq.Ack(pjst); err != nil {
		t.Fatal(err)
	}

	//job, exists := tst.jobs.Get()
	job, exists := tst.jobs.PollMT(time.Second * 10)
	if !exists || job == nil {
		t.Fatalf("expected received job")
	}

	if job.UID() != jst.UID() {
		t.Fatalf("actual UID %s expected UID %s", job.UID(), jst.UID())
	}

	sjs, err = tst.jq.Check(sjs)
	if err != nil {
		t.Fatal(err)
	}

	if sjs.State() != jobnik.InProcess {
		t.Fatalf("wrong in process job status %s", sjs.State().String())
	}

	if job.Handler() != jord.Handler() {
		t.Fatalf("different job handler")
	}

	if !reflect.DeepEqual(job.Attributes(), jord.Attributes()) {
		t.Fatalf("different job attributes")
	}

	if !bytes.Equal(job.Payload(), jord.Payload()) {
		t.Fatalf("different payload")
	}

}
