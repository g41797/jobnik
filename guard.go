// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import (
	"context"
	"fmt"
)

var _ Jobnik = &guard{}

type guard struct {
}

func (grd *guard) InitOnce(jsc string) error { return nil }

func (grd *guard) Process(cncl context.Context, job Job) (JobStatus, error) {
	return JobStatus{}, fmt.Errorf("Process is not implemented yet")
}

func (grd *guard) FinishOnce() error { return nil }
