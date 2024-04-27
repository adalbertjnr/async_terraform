package main

import (
	"log"
)

func main() {
	params := NewInputParser()
	if errors := params.Validate(); len(errors) > 0 {
		JSON(errors)
		return
	}

	svc, err := NewTerraformService(terraformInstaller)
	if err != nil {
		log.Fatal(err)
	}

	task := NewTaskManager(params, svc).retrieveTasks()
	task.start()
}
