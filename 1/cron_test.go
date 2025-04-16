package cron

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

type taskStub struct{}

func (taskStub) Exec() {
	log.Info().Msg("Executing Task")
}

func Test_Cron(t *testing.T) {
	for range 100 {
		Add(taskStub{}, time.Now().Add(time.Second))
	}
	time.Sleep(time.Second * 2)

	mux.Lock()
	if len(tasks) != 0 {
		t.Error("len(tasks) != 0")
	}
	mux.Unlock()
}
