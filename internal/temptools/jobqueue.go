package temptools

import (
	"errors"
	"log/slog"
	"time"
)

type CallbackFunc = func(id string, isRunning bool, lastRunTime time.Time)

type JobQueue struct {
	logger        slog.Logger
	recurringJobs map[string]*jobContainer
	c             chan *jobContainer
}

type Job interface {
	ID() string
	Run() error
	RunEvery() *time.Duration
}

type jobContainer struct {
	job         Job
	period      *time.Duration
	lastRunTime *time.Time
	callback    CallbackFunc
	isPushed    bool
}

func NewJobQueue(logger slog.Logger, size int) JobQueue {
	jobQueue := JobQueue{
		logger:        logger,
		recurringJobs: make(map[string]*jobContainer),
		c:             make(chan *jobContainer, size),
	}

	jobQueue.startWorker()
	jobQueue.startScheduler()

	return jobQueue
}

func (q *JobQueue) Push(job Job, callback CallbackFunc) error {
	jobContainer := &jobContainer{
		job:      job,
		period:   job.RunEvery(),
		callback: callback,
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

func (q *JobQueue) FetchState(id string) (bool, *time.Time) {
	rj, ok := q.recurringJobs[id]
	if !ok {
		return false, nil
	}

	return rj.isPushed, rj.lastRunTime
}

func (q *JobQueue) push(jobContainer *jobContainer) {
	jobContainer.isPushed = true
	q.c <- jobContainer
}

func (q *JobQueue) startWorker() {
	go func() {
		for {
			jobContainer := <-q.c
			err := jobContainer.run()
			if err != nil {
				q.logger.Error(err.Error())
			}
		}
	}()
}

func (q *JobQueue) startScheduler() {
	go func() {
		for {
			for k := range q.recurringJobs {
				job := q.recurringJobs[k]
				if job.shouldRun() {
					q.push(job)
				}
			}

			time.Sleep(getSmallestPeriod(q.recurringJobs))
		}
	}()
}

func getSmallestPeriod(jobContainers map[string]*jobContainer) time.Duration {
	var smallestPeriod *time.Duration = nil

	for _, c := range jobContainers {
		if smallestPeriod == nil || c.period.Nanoseconds() < smallestPeriod.Nanoseconds() {
			smallestPeriod = c.period
		}
	}

	if smallestPeriod == nil {
		return 10 * time.Second
	}

	return *smallestPeriod
}

func (c *jobContainer) run() error {
	defer func() {
		c.isPushed = false
	}()

	c.callback(c.job.ID(), true, *c.lastRunTime)

	nowUTC := time.Now().UTC()
	c.lastRunTime = &nowUTC

	err := c.job.Run()
	if err != nil {
		return err
	}

	c.callback(c.job.ID(), false, *c.lastRunTime)
	return nil
}

func (c jobContainer) shouldRun() bool {
	return !c.isPushed && (c.lastRunTime == nil || c.lastRunTime.Add(*c.period).After(time.Now().UTC()))
}
