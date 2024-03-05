// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik

import "github.com/g41797/sputnik"

type DefaultJobOrder struct {
	Name  string
	Atrbs []JobAttribute
	Pld   []byte
}

func (jo *DefaultJobOrder) Handler() string { return jo.Name }

func (jo *DefaultJobOrder) Attributes() []JobAttribute { return jo.Atrbs }

func (jo *DefaultJobOrder) Payload() []byte { return jo.Pld }

func (jo *DefaultJobOrder) OrderTo(m sputnik.Msg) error { return nil }

func (jo *DefaultJobOrder) OrderFrom(m sputnik.Msg) error { return nil }

type DefaultJobStatus struct {
	Uid string
	St  JobState
	Mem any
}

func (js *DefaultJobStatus) UID() string { return js.Uid }

func (js *DefaultJobStatus) State() JobState { return js.St }

func (js *DefaultJobStatus) Memento() any { return js.Mem }

func (jo *DefaultJobStatus) StatusTo(m sputnik.Msg) error { return nil }

func (jo *DefaultJobStatus) StatusFrom(m sputnik.Msg) error { return nil }

type DefaultJob struct {
	DefaultJobStatus
	DefaultJobOrder
}

var _ Job = &DefaultJob{}
