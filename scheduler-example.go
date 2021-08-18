package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GlobalLog string `yaml:"GlobalLog"`
	Jobs      []Job  `yaml:"Jobs"`
}

type Job struct {
	Name      string `yaml:"Name"`
	Command   string `yaml:"Command"`
	Frequency string `yaml:"Frequency"`
	Tag       string `yaml:"Tag,omitempty"`
	Log       string `yaml:"Log,omitempty"`
}

const ROOT string = "/Users/samwood/Documents/Code/Go/scheduler-example"

func main() {
	// STEP1: Parse the config file
	configFile := filepath.Join(ROOT, "config", "scheduler-example.yaml")
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	configData, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	c := Config{}
	err = yaml.Unmarshal(configData, &c)
	if err != nil {
		log.Fatalln(err)
	}

	// STEP2: Validate the config
	if c.GlobalLog == "" {
		log.Fatalln("Global log not defined")
	}
	if c.Jobs == nil {
		log.Fatalln("No jobs defined")
	}
	for _, j := range c.Jobs {
		if j.Name == "" {
			log.Fatalln("Job name not defined")
		}
		if j.Command == "" {
			log.Fatalln("Job command not defined")
		}
		if j.Frequency == "" {
			log.Fatalln("Job frequency not defined")
		}
	}

	logfile, err := os.OpenFile(c.GlobalLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	// STEP3: Create the scheduler and add jobs
	log.Println("Creating a scheduler")
	s := gocron.NewScheduler(time.UTC)
	log.Println("Created:", s)

	for _, j := range c.Jobs {
		script := filepath.Join(ROOT, "scripts", j.Command)
		log.Println("Script to execute:", script)
		cmd := exec.Command("/bin/bash", script)
		log.Println("Command:", cmd.String())
		cmdLog := filepath.Join(ROOT, "logs", j.Log)
		log.Println("Command log:", cmdLog)
		l, err := os.OpenFile(cmdLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		cmd.Stdout = l
		cmd.Stderr = l
		freq := j.Frequency
		tag := j.Tag
		log.Println("Adding job...")
		s.Every(freq).Tag(tag).Do(func() {
			err = cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
		})
		log.Println("Job added:", cmd.String(), tag, freq)
	}

	// STEP4: Run the scheduler
	log.Println("Running the scheduler...")
	s.StartBlocking()
}
