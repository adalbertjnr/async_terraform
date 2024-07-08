package main

import (
	"context"
	"log"
)

func main() {
	params := NewInputParser()
	if errors := params.Validate(); len(errors) > 0 {
		logErrors(errors)
		return
	}

	svc, err := NewTerraformService(terraformInstaller, params.version)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	task := NewTaskManager(ctx, params, svc).retrieveTasks()
	task.start()
}
