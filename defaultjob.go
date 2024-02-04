// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

type defaultJobOrder struct {
	name    string
	atrbs   []JobAttribute
	payload string
}

func (jo *defaultJobOrder) Name() string { return jo.name }

func (jo *defaultJobOrder) Attributes() []JobAttribute { return jo.atrbs }

func (jo *defaultJobOrder) Payload() string { return jo.payload }

type defaultJob struct {
	id string
	defaultJobOrder
}

func (job *defaultJob) UID() string { return job.id }
