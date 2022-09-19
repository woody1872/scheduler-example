package schedulerexample

import (
	"errors"
	"io"
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	ErrMissingGlobalLog    = errors.New("missing global log file")
	ErrMissingJobCommand   = errors.New("missing job command")
	ErrMissingJobFrequency = errors.New("missing job frequency")
	ErrMissingJobType      = errors.New("missing job type")
)

var (
	ValidFrequency = regexp.MustCompile("^[1-9][0-9]*[sm]$")
)

// Config represents the scheduler config file, where jobs and their
// required properties get defined to be run later by the scheduler
type Config struct {
	GlobalLog string `yaml:"GlobalLog"`
	Jobs      []Job  `yaml:"Jobs"`
}

// ReadConfig reads a config file and unmarshals it in to the Config type
// which represents a scheduler config
func ReadConfig(cf io.Reader) (Config, error) {
	data, err := io.ReadAll(cf)
	if err != nil {
		log.Fatalln(err)
	}

	c := Config{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// Valid is a helper function which checks whether a Config object, and subquently
// the underlying config file, is valid with all the required properties
func Valid(c *Config) bool {
	var v bool = true

	if v = hasCommand(c); !v {
		return v
	}
	if v = hasFrequency(c); !v {
		return v
	}
	if v = hasType(c); !v {
		return v
	}

	return v
}

func hasCommand(c *Config) bool {
	for _, job := range c.Jobs {
		if job.Command == "" {
			return false
		}
	}

	return true
}

func hasFrequency(c *Config) bool {
	for _, job := range c.Jobs {
		if job.Frequency == "" {
			return false
		}
		if !ValidFrequency.MatchString(job.Frequency) {
			return false
		}
	}

	return true
}

func hasType(c *Config) bool {
	for _, job := range c.Jobs {
		if job.Type == "" {
			return false
		}
	}

	return true
}
