package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	schedulerexample "github.com/SWOOD1872/scheduler-example"
	"github.com/go-co-op/gocron"
)

const ROOT string = "/Users/samwood/Documents/Code/Go/scheduler-example"

func main() {
	rootCfg := filepath.Join(ROOT, "config", "scheduler-example.yaml")
	f, err := os.Open(rootCfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	cfg, err := schedulerexample.ReadConfig(f)
	if err != nil {
		log.Fatalln(err)
	}

	s := gocron.NewScheduler(time.UTC)
	err = schedulerexample.AddJobs(s, cfg.Jobs)
	if err != nil {
		log.Fatalln(err)
	}
}
