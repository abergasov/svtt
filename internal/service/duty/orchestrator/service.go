package orchestrator

import (
	"context"
	"svtt/internal/entities"
	"svtt/internal/logger"
	"svtt/internal/service/duty/processor"
	"sync"

	"github.com/google/uuid"
)

type ProcessorCreationFunc func(ctx context.Context, log logger.AppLogger, responseChan chan entities.DutyResponce, processor processor.SingleDutyProcessor) processor.Dutyer

type Service struct {
	ctx             context.Context
	cancel          context.CancelFunc
	log             logger.AppLogger
	requestChan     chan entities.DutyRequest
	responseChan    chan entities.DutyResponce
	dutyProcessor   map[int]processor.Dutyer
	dutyProcessorMU sync.RWMutex
	creator         ProcessorCreationFunc
	singleProcessor processor.SingleDutyProcessor
}

type Option func(*Service)

func WithCustomRequestProcessor(custom processor.SingleDutyProcessor) func(*Service) {
	return func(s *Service) {
		s.singleProcessor = custom
	}
}

func NewService(log logger.AppLogger, creator ProcessorCreationFunc, opts ...Option) *Service {
	srv := &Service{
		log:           log,
		requestChan:   make(chan entities.DutyRequest, 1_000),
		responseChan:  make(chan entities.DutyResponce, 1_000),
		dutyProcessor: make(map[int]processor.Dutyer),
		creator:       creator,
		singleProcessor: func(request entities.DutyRequest) entities.DutyResponce {
			return entities.DutyResponce{
				Duty:      request.Duty,
				Height:    request.Height,
				Validator: request.Validator,
				Response:  uuid.NewString(),
			}
		},
	}
	for _, opt := range opts {
		opt(srv)
	}
	srv.ctx, srv.cancel = context.WithCancel(context.Background())
	go srv.processBackground()
	return srv
}

func (s *Service) HandleDutyRequest(request entities.DutyRequest) {
	s.requestChan <- request
}

func (s *Service) GetDutyResponser() chan entities.DutyResponce {
	return s.responseChan
}

func (s *Service) Stop() error {
	s.cancel()
	s.dutyProcessorMU.RLock()
	defer s.dutyProcessorMU.RUnlock()
	for _, validatorProcessor := range s.dutyProcessor {
		validatorProcessor.Stop()
	}
	return nil
}
