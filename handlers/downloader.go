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
}

func (dl *downloader) InitOnce(_ string) error {
	return nil
}

func (dl *downloader) FinishOnce() error {
	return nil
}

func (dl *downloader) Process(cncl context.Context, job jobnik.Job) (jobnik.JobStatus, error) {
	return jobnik.JobStatus{}, fmt.Errorf("Process for %s is not implemented", downloaderJobnikName)
}
