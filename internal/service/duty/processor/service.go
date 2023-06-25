package processor

import (
	"container/list"
	"context"
	"svtt/internal/entities"
	"svtt/internal/logger"
	"sync"
	"time"
)

type Service struct {
	wg               sync.WaitGroup
	ticker           *time.Ticker
	ctx              context.Context
	log              logger.AppLogger
	responseChan     chan entities.DutyResponce
	executionQueue   map[entities.DutyType]*list.List
	executionQueueMU sync.RWMutex
	processor        SingleDutyProcessor
}

func NewService(ctx context.Context, log logger.AppLogger, responseChan chan entities.DutyResponce, processor SingleDutyProcessor) Dutyer {
	srv := &Service{
		ctx:            ctx,
		log:            log,
		responseChan:   responseChan,
		executionQueue: make(map[entities.DutyType]*list.List),
		ticker:         time.NewTicker(100 * time.Millisecond),
		processor:      processor,
	}
	go srv.processBackground()
	return srv
}

func (s *Service) Process(request entities.DutyRequest) {
	s.executionQueueMU.Lock()
	defer s.executionQueueMU.Unlock()
	queue, ok := s.executionQueue[request.Duty]
	if !ok {
		queue = list.New()
		s.executionQueue[request.Duty] = queue
	}
	queue.PushBack(&processingContainer{
		status:  processingStatusReady,
		request: request,
	})
}

func (s *Service) Stop() {
	s.wg.Wait()
}
