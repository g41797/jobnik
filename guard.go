// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
)

var _ Jobnik = &guard{}

type guardState int

const (
	initallowed          guardState = iota
	processfinishallowed guardState = iota + 1
	nothingallowed       guardState = iota + 2
)

func (gst guardState) String() string {
	return []string{"initallowed", "processfinishallowed", "nothingallowed"}[gst]
}

type guard struct {
	lock  sync.Mutex
	state guardState
	name  string
	jbnk  Jobnik
}

func (grd *guard) tryCreate(name string) (Jobnik, error) {

	grd.name = name

	fact, exists := factories.Load(strings.ToLower(name))
	if !exists {
		return nil, fmt.Errorf("factory for %s does not exist", grd.name)
	}

	if jbnk, err := fact.(JobnikFactory)(); err == nil {
		grd.jbnk = jbnk
		grd.state = initallowed
		return grd, nil
	} else {
		return nil, err
	}
}

func (grd *guard) InitOnce(jsc string) error {
	if grd == nil {
		return fmt.Errorf("InitOnce nil guard")
	}

	grd.lock.Lock()
	defer grd.lock.Unlock()

	if grd.jbnk == nil {
		return fmt.Errorf("jobnik was not created")
	}

	if grd.state != initallowed {
		return fmt.Errorf("Init disabled for %s", grd.state.String())
	}

	if err := grd.jbnk.InitOnce(jsc); err != nil {
		grd.state = nothingallowed
		return err
	}

	grd.state = processfinishallowed
	return nil
}

func (grd *guard) Process(ctx context.Context, job Job) (JobStatus, error) {
	if grd == nil {
		return JobStatus{}, fmt.Errorf("Process nil guard")
	}

	grd.lock.Lock()
	defer grd.lock.Unlock()

	if grd.jbnk == nil {
		return JobStatus{}, fmt.Errorf("jobnik was not created")
	}

	if grd.state != processfinishallowed {
		return JobStatus{}, fmt.Errorf("Process disabled for %s", grd.state.String())
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return grd.jbnk.Process(ctx, job)
}

func (grd *guard) FinishOnce() error {
	if grd == nil {
		return nil
	}

	grd.lock.Lock()
	defer grd.lock.Unlock()

	if grd.jbnk == nil {
		return nil
	}

	if grd.state == initallowed {
		return nil
	}

	if grd.state != processfinishallowed {
		return fmt.Errorf("Finish disabled for %s", grd.state.String())
	}

	err := grd.jbnk.FinishOnce()

	grd.state = nothingallowed

	return err
}

func storeFactory(name string, fact JobnikFactory) {
	if len(name) == 0 {
		log.Panic("empty jobnik name")
	}
	if fact == nil {
		log.Panicf("nil jobnik factory for %s", name)
	}

	if _, exists := factories.LoadOrStore(strings.ToLower(name), fact); exists {
		log.Panicf("jobnik factory for %s already exists", name)
	}
	return
}

var factories sync.Map
