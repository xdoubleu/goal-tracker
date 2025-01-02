package services

import (
	"context"
	"net/http"
	"strconv"
	"time"

	wstools "github.com/XDoubleU/essentia/pkg/communication/ws"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/temptools"
)

type WebSocketService struct {
	allowedOrigins []string
	handler        *wstools.WebSocketHandler[dtos.SubscribeMessageDto]
	jobQueue       *temptools.JobQueue
	topics         map[string]*wstools.Topic
}

func NewWebSocketService(
	allowedOrigins []string,
	jobQueue *temptools.JobQueue,
) *WebSocketService {
	service := WebSocketService{
		allowedOrigins: allowedOrigins,
		handler:        nil,
		jobQueue:       jobQueue,
		topics:         make(map[string]*wstools.Topic),
	}

	handler := wstools.CreateWebSocketHandler[dtos.SubscribeMessageDto](
		1,
		100, //nolint:mnd //no magic number
	)

	service.handler = &handler
	service.registerTopics()

	return &service
}

func (service WebSocketService) Handler() http.HandlerFunc {
	return service.handler.Handler()
}

func (service WebSocketService) UpdateState(
	id string,
	isRunning bool,
	lastRunTime *time.Time,
) {
	topic, ok := service.topics[id]
	if !ok {
		return
	}

	topic.EnqueueEvent(dtos.StateMessageDto{
		IsRefreshing: isRunning,
		LastRefresh:  lastRunTime,
	})
}

func (service WebSocketService) registerTopics() {
	topics := []string{
		"todoist",
		strconv.Itoa(int(models.SteamCompletionRate.ID)),
	}

	for _, topic := range topics {
		registeredTopic, err := service.handler.AddTopic(
			topic,
			service.allowedOrigins,
			func(_ context.Context, tp *wstools.Topic) (any, error) {
				return service.fetchState(tp), nil
			},
		)
		if err != nil {
			panic(err)
		}
		service.topics[topic] = registeredTopic
	}
}

func (service WebSocketService) fetchState(topic *wstools.Topic) dtos.StateMessageDto {
	isRefreshing, lastRefresh := service.jobQueue.FetchState(topic.Name)

	return dtos.StateMessageDto{
		IsRefreshing: isRefreshing,
		LastRefresh:  lastRefresh,
	}
}
