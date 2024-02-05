package main

import (
	"embed"
	"fmt"

	"github.com/g41797/jobnik"
	_ "github.com/g41797/jobnik/handlers"
)

//go:embed _configs/*.json
var configs embed.FS

func main() {
	dlr, _ := jobnik.NewJobnik("DownloadFiles")
	defer dlr.FinishOnce()

	cnfstr, err := configs.ReadFile("_configs/server.json")
	if err != nil {
		fmt.Print(err)
		return
	}

	err = dlr.InitOnce(string(cnfstr))
	if err != nil {
		fmt.Print(err)
		return
	}

}
