package temptools

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/XDoubleU/essentia/pkg/sentry"
)

type CallbackFunc = func(id string, isRunning bool, lastRunTime *time.Time)

type JobQueue struct {
	logger              slog.Logger
	recurringJobs       map[string]*jobContainer
	c                   chan *jobContainer
	workerStopRequested chan bool
	workerActive        bool
	schedulerActive     bool
}

type Job interface {
	ID() string
	Run(slog.Logger) error
	RunEvery() *time.Duration
}

type jobContainer struct {
	job         Job
	period      *time.Duration
	lastRunTime *time.Time
	callback    CallbackFunc
	isPushed    bool
}

func NewJobQueue(logger slog.Logger, size int) *JobQueue {
	jobQueue := &JobQueue{
		logger:              logger,
		recurringJobs:       make(map[string]*jobContainer),
		c:                   make(chan *jobContainer, size),
		workerStopRequested: make(chan bool),
		workerActive:        true,
		schedulerActive:     true,
	}

	jobQueue.startWorker()
	jobQueue.startScheduler()

	return jobQueue
}

func (q *JobQueue) Clear() {
	q.schedulerActive = false
	q.workerStopRequested <- true
	q.recurringJobs = make(map[string]*jobContainer)

	for q.workerActive {
		time.Sleep(100 * time.Millisecond)
	}
}

func (q *JobQueue) Push(job Job, callback CallbackFunc) error {
	jobContainer := &jobContainer{
		job:         job,
		period:      job.RunEvery(),
		callback:    callback,
		lastRunTime: nil,
		isPushed:    false,
	}

	if jobContainer.period != nil {
		_, ok := q.recurringJobs[job.ID()]
		if ok {
			return errors.New("a job with this ID already exists")
		}

		q.recurringJobs[job.ID()] = jobContainer
	}

	q.push(jobContainer)
	return nil
}

func (q *JobQueue) ForceRun(id string) {
	rj, ok := q.recurringJobs[id]
	if !ok {
		return
	}
	q.push(rj)
}

func (q *JobQueue) FetchRecurringJobIDs() []string {
	result := []string{}
	for _, rj := range q.recurringJobs {
		result = append(result, rj.job.ID())
	}
	return result
}

func (q *JobQueue) FetchState(id string) (bool, *time.Time) {
	rj, ok := q.recurringJobs[id]
	if !ok {
		return false, nil
	}

	return rj.isPushed, rj.lastRunTime
}

func (q *JobQueue) push(jobContainer *jobContainer) {
	if !q.schedulerActive {
		q.startScheduler()
	}

	if !q.workerActive {
		q.startWorker()
	}

	jobContainer.isPushed = true
	q.c <- jobContainer
}

func (q *JobQueue) startWorker() {
	q.workerActive = true

	go sentry.GoRoutineErrorHandler(
		context.Background(),
		"JobQueueWorker",
		func(_ context.Context) error {
		out:
			for {
				select {
				case jobContainer := <-q.c:
					err := jobContainer.run(q.logger)
					if err != nil {
						q.logger.Error(err.Error())
					}
				case <-q.workerStopRequested:
					break out
				}
			}

			q.workerActive = false
			return nil
		},
	)
}

func (q *JobQueue) startScheduler() {
	q.schedulerActive = true

	go sentry.GoRoutineErrorHandler(
		context.Background(),
		"JobQueueScheduler",
		func(_ context.Context) error {
			for q.schedulerActive {
				for k := range q.recurringJobs {
					job := q.recurringJobs[k]
					if job.shouldRun() {
						q.push(job)
					}
				}
				time.Sleep(getSmallestPeriod(q.recurringJobs))
			}
			return nil
		},
	)
}

func getSmallestPeriod(jobContainers map[string]*jobContainer) time.Duration {
	var smallestPeriod *time.Duration

	for _, c := range jobContainers {
		if smallestPeriod == nil ||
			c.period.Nanoseconds() < smallestPeriod.Nanoseconds() {
			smallestPeriod = c.period
		}
	}

	if smallestPeriod == nil {
		//nolint:mnd //no magic number
		return 10 * time.Second
	}

	return *smallestPeriod
}

func (c *jobContainer) run(logger slog.Logger) error {
	defer func() {
		c.isPushed = false
	}()

	c.callback(c.job.ID(), true, c.lastRunTime)

	nowUTC := time.Now().UTC()
	c.lastRunTime = &nowUTC

	logger.Debug(fmt.Sprintf("started job %s", c.job.ID()))
	err := c.job.Run(logger)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("successfully finished job %s", c.job.ID()))

	c.callback(c.job.ID(), false, c.lastRunTime)
	return nil
}

func (c jobContainer) shouldRun() bool {
	return !c.isPushed &&
		(c.lastRunTime == nil || c.lastRunTime.Add(*c.period).After(time.Now().UTC()))
}
