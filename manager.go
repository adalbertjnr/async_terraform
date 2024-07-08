package main

import (
	"context"
	"log"
	"log/slog"
	"strings"
)

const (
	Apply   = "apply"
	Destroy = "destroy"
	Plan    = "plan"
)

type TaskManager struct {
	ctx        context.Context
	tfsvc      TFService
	dataParams InputParser
	tasks      []string
	taskch     chan string
	done       chan struct{}
}

func NewTaskManager(ctx context.Context, input InputParser, svc *TFService) *TaskManager {
	return &TaskManager{
		ctx:        ctx,
		dataParams: input,
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
	slog.Info("message", "started by user", w.dataParams.user, "verb", strings.ToUpper(w.dataParams.verb))

	w.startChannels()
	w.bootstrap(w.dataParams.workers)

	close(w.taskch)
	w.release(w.dataParams.workers)
}

func (w *TaskManager) bootstrap(workers int) {
	for i := 0; i < workers; i++ {
		go w.worker(i, w.dataParams.verb)
	}

	for _, task := range w.tasks {
		w.taskch <- task
	}
}

func (w *TaskManager) release(workers int) {
	for i := 0; i < workers; i++ {
		<-w.done
	}
}

func (t *TaskManager) worker(workerId int, verb string) {
	defer func() {
		t.done <- struct{}{}
	}()

	for task := range t.taskch {
		slog.Info("running terraform", "aws account", task, "worker", workerId, "verb", strings.ToUpper(verb))

		switch verb {
		case Apply:
			err := t.tfsvc.terraformTaskCreate(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		case Destroy:
			err := t.tfsvc.terraformTaskDestroy(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		case Plan:
			err := t.tfsvc.terraformTaskPlan(t.ctx, task, t.tfsvc.execPath)
			if err != nil {
				log.Println(err)
			}
		default:
			slog.Warn("task manager warn", "verb", strings.ToUpper(verb), "status", "not found")
		}
	}
}

func (w *TaskManager) startChannels() {
	w.taskch = make(chan string, len(w.tasks))
	w.done = make(chan struct{})
}
