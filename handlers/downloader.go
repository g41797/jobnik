package handlers

import (
	"context"
	"fmt"

	"github.com/g41797/jobnik"
)

const downloaderJobnikName = "DownloadFiles"

// Factory function
func downloaderFactory() (jobnik.Jobnik, error) {
	return new(downloader), nil
}

// Register handler:
func init() {
	jobnik.RegisterFactory(downloaderJobnikName, downloaderFactory)
}

type downloader struct {
	conf serverConfig
}

func (dl *downloader) InitOnce(js string) error {
	return dl.conf.unmarshall(js)
}

func (dl *downloader) FinishOnce() error {
	return nil
}

func (dl *downloader) Process(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	return jobnik.JobStatus{}, fmt.Errorf("Process for %s is not implemented", downloaderJobnikName)
}
