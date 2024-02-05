package main

import (
	"fmt"

	"github.com/g41797/jobnik"
	_ "github.com/g41797/jobnik/handlers"
)

func main() {
	worker, _ := jobnik.NewJobnik("DownloadFiles")
	defer worker.FinishOnce()
	fmt.Printf("Type of worker %T", worker)
}
