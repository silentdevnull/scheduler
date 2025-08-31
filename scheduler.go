// Package scheduler  to run go function or other thing as needed
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Schedule interface {
	Next(t time.Time) time.Time
}

type Job struct {
	schedule Schedule
	task     func()
	nextRun  time.Time
}

type Scheduler struct {
	jobs    []*Job
	mu      sync.Mutex
	stop    chan struct{}
	running bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make([]*Job, 0),
		stop: make(chan struct{}),
	}
}

func (s *Scheduler) AddJob(schedule Schedule, task func()) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	job := &Job{
		schedule: schedule,
		task:     task,
		nextRun:  schedule.Next(time.Now()),
	}
	s.jobs = append(s.jobs, job)
	return job
}

func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()
	go s.run()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	close(s.stop)
	s.running = false
	s.mu.Unlock()
}

func (s *Scheduler) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			s.runPending(now)
		case <-s.stop:
			return
		}
	}
}

func (s *Scheduler) runPending(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, job := range s.jobs {
		if job.nextRun.Before(now) || job.nextRun.Equal(now) {
			go job.task()
			job.nextRun = job.schedule.Next(job.nextRun)
		}
	}
}

type intervalSchedule struct {
	interval time.Duration
}

func (is intervalSchedule) Next(t time.Time) time.Time {
	return t.Add(is.interval)
}

func Daily(hour, min, sec int) Schedule {
	return &dailySchedule{
		hour: hour,
		min:  min,
		sec:  sec,
	}
}

type dailySchedule struct {
	hour, min, sec int
}

func (ds dailySchedule) Next(from time.Time) time.Time {
	now := time.Now()
	next := time.Date(
		from.Year(),
		from.Month(),
		from.Day(),
		ds.hour,
		ds.min,
		ds.sec,
		0,
		now.Location(),
	)
	if next.After(from) {
		return next
	}
	return next.AddDate(0, 0, 1)
}

func AtHour(n int) Schedule {
	if n <= 0 {
		n = 1
	}
	return &atHourSchedule{n: n}
}

type atHourSchedule struct {
	n int
}

func (ahs atHourSchedule) Next(from time.Time) time.Time {
	next := from.Truncate(time.Hour)
	for !next.After(from) {
		next = next.Add(time.Duration(ahs.n) * time.Hour)
	}
	return next
}

func AtMinute(n int) Schedule {
	if n <= 0 {
		n = 1
	}
	return intervalSchedule{interval: time.Duration(n) + time.Minute}
}

func Cron(expression string) (Schedule, error) {
	fields := strings.Fields(expression)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	cs := &cronSchedule{}
	var err error
	cs.minute, err = parseCronField(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}
	cs.hour, err = parseCronField(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}
	if fields[2] != "*" || fields[3] != "*" || fields[4] != "*" {
		return nil, fmt.Errorf("only '*' is supported for day, month, and day-of-week fields in this simplified version")
	}

	return cs, nil
}

type cronSchedule struct {
	minute, hour int
}

func (cs cronSchedule) Next(from time.Time) time.Time {
	next := from.Add(time.Minute).Truncate(time.Minute)

	for {
		if cs.hour != -1 && next.Hour() != cs.hour {
			next = next.Truncate(time.Hour).Add(time.Hour)
			continue
		}

		if cs.minute != -1 && next.Minute() != cs.minute {
			next = next.Add(time.Minute)
			continue
		}
		return next
	}
}

func parseCronField(field string) (int, error) {
	if field == "*" {
		return -1, nil
	}
	val, err := strconv.Atoi(field)
	if err != nil {
		return 0, err
	}
	return val, nil
}
