package temptools_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/temptools"
)

type TestJob struct {
	isRecurring bool
}

func (j TestJob) ID() string {
	return "test"
}

func (j TestJob) Run(_ *slog.Logger) error {
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (j TestJob) RunEvery() *time.Duration {
	delay := 500 * time.Millisecond

	if j.isRecurring {
		return &delay
	}

	return nil
}

func TestJobQueueSimple(t *testing.T) {
	jobQueue := temptools.NewJobQueue(logging.NewNopLogger(), 1)

	states := []bool{}

	err := jobQueue.Push(
		TestJob{isRecurring: false},
		func(_ string, isRunning bool, _ *time.Time) {
			states = append(states, isRunning)
		},
	)
	assert.Nil(t, err)

	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, []bool{true, false}, states)
}

func TestJobQueueSimpleAfterClear(t *testing.T) {
	jobQueue := temptools.NewJobQueue(logging.NewNopLogger(), 1)

	states := []bool{}

	err := jobQueue.Push(
		TestJob{isRecurring: false},
		func(_ string, isRunning bool, _ *time.Time) {
			states = append(states, isRunning)
		},
	)
	assert.Nil(t, err)

	time.Sleep(1 * time.Millisecond)
	jobQueue.Clear()
	assert.Equal(t, []bool{true, false}, states)

	err = jobQueue.Push(
		TestJob{isRecurring: false},
		func(_ string, isRunning bool, _ *time.Time) {
			states = append(states, isRunning)
		},
	)
	assert.Nil(t, err)

	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, []bool{true, false, true, false}, states)
}

func TestJobQueueRecurring(t *testing.T) {
	jobQueue := temptools.NewJobQueue(logging.NewNopLogger(), 1)

	states := []bool{}

	err := jobQueue.Push(
		TestJob{isRecurring: true},
		func(_ string, isRunning bool, _ *time.Time) {
			states = append(states, isRunning)
		},
	)
	assert.Nil(t, err)

	jobIDs := jobQueue.FetchRecurringJobIDs()
	assert.Equal(t, []string{"test"}, jobIDs)

	state, _ := jobQueue.FetchState("test")
	assert.Equal(t, true, state)

	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, []bool{true, false}, states)

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, []bool{true, false, true, false}, states)
}

func TestJobQueueRecurringForce(t *testing.T) {
	jobQueue := temptools.NewJobQueue(logging.NewNopLogger(), 1)

	states := []bool{}

	err := jobQueue.Push(
		TestJob{isRecurring: true},
		func(_ string, isRunning bool, _ *time.Time) {
			states = append(states, isRunning)
		},
	)
	assert.Nil(t, err)

	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, []bool{true, false}, states)

	jobQueue.ForceRun("test")

	state, _ := jobQueue.FetchState("test")
	assert.Equal(t, true, state)

	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, []bool{true, false, true, false}, states)
}
