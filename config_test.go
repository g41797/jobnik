// Copyright (c) 2024 g41797
// SPDX-License-Identifier: MIT

package jobnik_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

//-------------------------------------------------------------------------------------------------------
// Useful JSON&Go links:
//
// https://jsonlint.com
// https://betterstack.com/community/guides/scaling-go/json-in-go/
// https://www.digitalocean.com/community/tutorials/how-to-use-json-in-go#using-a-struct-to-generate-json
// https://blog.boot.dev/golang/json-golang/
//-------------------------------------------------------------------------------------------------------

// bool needs special attention:
// see https://biscuit.ninja/posts/unmarshalling-json-with-null-booleans-in-go/#solution-1-use-pointers
type config struct {
	IsLogErrors   *bool `json:"errors,omitempty"`
	IsLogWarnings *bool `json:"warnings,omitempty"`
	IsLogInfo     *bool `json:"info,omitempty"`
}

func (cnf *config) setDefault() {
	cnf.IsLogErrors = new(bool)
	*cnf.IsLogErrors = true

	cnf.IsLogWarnings = nil
	cnf.IsLogInfo = nil
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

func TestConfig(t *testing.T) {
	var cnf config

	if err := cnf.unmarshall("wrong json"); err == nil {
		t.Errorf("unmarshall should return error for wrong JSON")
	}

	onlyerrfalse := `{
		"errors": false
	}`

	if err := cnf.unmarshall(onlyerrfalse); err != nil {
		t.Errorf("unmarshall should not return error for valid JSON")
	}

	if cnf.isLogWarnings() || cnf.isLogInfo() || cnf.isLogErrors() {
		t.Errorf("wrong onlyerrfalse unmarshallibg")
	}

	onlyerrtrue := `{
		"errors": true
	}`

	if err := cnf.unmarshall(onlyerrtrue); err != nil {
		t.Errorf("unmarshall should not return error for valid JSON")
	}

	if cnf.isLogWarnings() || cnf.isLogInfo() || !cnf.isLogErrors() {
		t.Errorf("wrong onlyerrtrue unmarshallibg")
	}

	alltrue := `{
		"errors": true,
		"warnings": true,
		"info":true
	}`

	if err := cnf.unmarshall(alltrue); err != nil {
		t.Errorf("unmarshall should not return error for valid JSON")
	}

	if !cnf.isLogWarnings() || !cnf.isLogInfo() || !cnf.isLogErrors() {
		t.Errorf("wrong alltrue unmarshallibg")
	}
}
