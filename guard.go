// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var _ Jobnik = &guard{}

type guard struct {
	lock sync.Mutex
	name string
	jbnk Jobnik
}

func (grd *guard) InitOnce(jsc string) error { return nil }

func (grd *guard) Process(cncl context.Context, job Job) (JobStatus, error) {
	return JobStatus{}, fmt.Errorf("Process is not implemented yet")
}

func (grd *guard) FinishOnce() error { return nil }

func (grd *guard) tryCreate(name string) (Jobnik, error) {

	grd.name = name

	fact, exists := factories.Load(strings.ToLower(name))
	if !exists {
		return nil, fmt.Errorf("factor for %s does not exist", grd.name)
	}

	if jbnk, err := fact.(JobnikFactory)(); err == nil {
		grd.jbnk = jbnk
		return grd, nil
	} else {
		return nil, err
	}
}

func storeFactory(name string, fact JobnikFactory) {
	if len(name) == 0 {
		panic("empty jobnik name")
	}
	if fact == nil {
		panic(fmt.Errorf("nil jobnik factory for %s", name))
	}

	if _, exists := factories.LoadOrStore(strings.ToLower(name), fact); exists {
		panic(fmt.Errorf("jobnik factory for %s already exists", name))
	}
	return
}

var factories sync.Map
