package scheduler

import "testing"

func TestSchedulerStartStop(t *testing.T) {
	s := NewScheduler()
	s.Start()
	if !s.running {
		t.Error("Scheduler should be running after Start()")
	}
	s.Stop()
	if s.running {
		t.Error("Scheduler should be stopped after Stop()")
	}
}
