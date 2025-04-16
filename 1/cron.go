package cron

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	mux   sync.Mutex
	tasks map[Task]struct{} = make(map[Task]struct{})
)

type Task interface {
	Exec()
}

func init() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()
}

func Add(task Task, t time.Time) {
	if time.Now().After(t) {
		return
	}

	mux.Lock()
	defer mux.Unlock()
	tasks[task] = struct{}{}
	go runTask(task, t)
}

// runTask ожидает пока не наступит время выполнения задания,
// запускает его выполнение и удаляет из списка запланированных заданий.
func runTask(task Task, startTime time.Time) {
	time.Sleep(time.Until(startTime))

	task.Exec()

	mux.Lock()
	delete(tasks, task)
	mux.Unlock()
}
