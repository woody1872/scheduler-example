package schedulerexample

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/go-co-op/gocron"
)

var (
	ErrUnknownJobType = errors.New("unknown job type")
)

type Job struct {
	Command   string `yaml:"Command"`
	Frequency string `yaml:"Frequency"`
	Type      string `yaml:"Type"`
	Tag       string `yaml:"Tag,omitempty"`
	Log       string `yaml:"Log,omitempty"`
}

func AddJob(sched *gocron.Scheduler, job Job) error {
	switch strings.ToLower(job.Type) {
	case "bash":
		sched.Every(job.Frequency).Tag(job.Tag).Do(RunBash, job)
	case "python":
		sched.Every(job.Frequency).Tag(job.Tag).Do(RunPython, job)
	default:
		return ErrUnknownJobType
	}

	return nil
}

func AddJobs(sched *gocron.Scheduler, jobs []Job) error {
	for _, job := range jobs {
		err := AddJob(sched, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunBash(j Job) error {
	jobLog, err := os.OpenFile(j.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/bash", j.Command)
	cmd.Stdout = jobLog
	cmd.Stderr = jobLog
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func RunPython(j Job) error {
	jobLog, err := os.OpenFile(j.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/bin/python3", j.Command)
	cmd.Stdout = jobLog
	cmd.Stderr = jobLog
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
