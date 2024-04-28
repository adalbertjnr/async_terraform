package main

import (
	"context"
	"log"
	"strings"
)

const (
	Apply   = "apply"
	Destroy = "destroy"
	Plan    = "plan"
)

type TaskManager struct {
	tfsvc      TFService
	dataParams workManagerParams
	tasks      []string
	taskch     chan string
	done       chan struct{}
}

type workManagerParams struct {
	user    string
	verb    string
	tasks   string
	version string
	workers int
}

func managerData(input InputParser) workManagerParams {
	return workManagerParams(input)
}

func NewTaskManager(input InputParser, svc *TFService) *TaskManager {
	return &TaskManager{
		dataParams: managerData(input),
		tfsvc:      *svc,
	}
}

func (w *TaskManager) retrieveTasks() *TaskManager {
	data := w.dataParams.tasks
	tasks := optionsGetter(data)

	if len(tasks) > 0 {
		w.tasks = optionsGetter(data)
		return w
	}

	panic(ErrEmptyTask)
}

func (w *TaskManager) start() {
	log.Printf("started by [%s] with verb [%s]\n", w.dataParams.user, strings.ToUpper(w.dataParams.verb))

	w.startChannels()
	workers := w.dataParams.workers
	for i := 0; i < workers; i++ {
		go w.worker(i, w.dataParams.verb)
	}

	for _, task := range w.tasks {
		w.taskch <- task
	}

	close(w.taskch)
	w.release(workers)
}

func (w *TaskManager) release(workers int) {
	for i := 0; i < workers; i++ {
		<-w.done
	}
}

func (t *TaskManager) worker(workerId int, verb string) {
	ctx := context.Background()
	for task := range t.taskch {
		log.Printf("running terraform on [%s] by worker [%d] with verb [%s]\n", task, workerId, strings.ToUpper(verb))

		switch verb {
		case Apply:
			err := t.tfsvc.terraformTaskCreate(ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		case Destroy:
			err := t.tfsvc.terraformTaskDestroy(ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		case Plan:
			err := t.tfsvc.terraformTaskPlan(ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("verb not found")
		}
	}
	t.done <- struct{}{}
}

func (w *TaskManager) startChannels() {
	w.taskch = make(chan string, len(w.tasks))
	w.done = make(chan struct{})
}
