package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
)

const (
	WorkflowPath        = ".github/workflows"
	ValidateTaskMessage = "task"
	ValidateVerbMessage = "verb"
)

type InputParser struct {
	user    string
	verb    string
	tasks   string
	version string
	workers int
}

func (t *InputParser) Validate() map[string]string {
	errors := make(map[string]string)
	if t.tasks == "" {
		errors[ValidateTaskMessage] = ErrEmptyTask.Error()
	}
	if !validateVerb(t.verb) {
		errors[ValidateVerbMessage] = ErrVerbNotFound.Error()
	}
	return errors
}

func JSON(data interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(data); err != nil {
		log.Println("error enconding the error message")
		return
	}
}

func validateVerb(inputVerb string) bool {
	verbs := map[string]bool{
		Plan:    true,
		Apply:   true,
		Destroy: true,
	}
	return verbs[strings.ToLower(inputVerb)]
}

func NewInputParser() InputParser {
	user := flag.String("user", "terraform", "the username who started the pipeline")
	verb := flag.String("verb", "plan", "set the verb - plan - apply - destroy")
	workers := flag.Int("workers", 2, "set the number of workers to run concurrently")
	tasks := flag.String("tasks", "", "tasks that will be executed by terraform")
	version := flag.String("version", "1.7.0", "terraform version")
	flag.Parse()
	return InputParser{
		user:    sanitizeInput(user),
		verb:    sanitizeInput(verb),
		tasks:   strings.TrimSpace(*tasks),
		version: strings.TrimSpace(*version),
		workers: *workers,
	}
}

func sanitizeInput(s *string) string {
	return strings.ToLower(strings.TrimSpace(*s))
}

func optionsGetter(data string) []string {
	var (
		tasks          = make([]string, 0)
		sanitizedTasks = make([]string, 0)
	)

	if len(data) > 0 {
		lines := strings.Split(data, "\n")
		for _, line := range lines {
			tasks = append(tasks, strings.Split(line, ",")...)
		}
		for _, task := range tasks {
			if task != "" {
				sanitizedTasks = append(sanitizedTasks, task)
			}
		}
	}
	return sanitizedTasks
}
