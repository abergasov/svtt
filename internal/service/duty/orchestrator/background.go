package orchestrator

import "go.uber.org/zap"

func (s *Service) processBackground() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case request := <-s.requestChan:
			s.log.Info("got request",
				zap.Int("height", request.Height),
				zap.Int("validator", request.Validator),
				zap.String("duty", request.Duty.String()),
			)
			s.dutyProcessorMU.Lock()
			validatorProcessor, ok := s.dutyProcessor[request.Validator]
			if !ok {
				validatorProcessor = s.creator(s.ctx, s.log.With(zap.Int("validator", request.Validator)), s.responseChan, s.singleProcessor)
				s.dutyProcessor[request.Validator] = validatorProcessor
			}
			s.dutyProcessorMU.Unlock()
			validatorProcessor.Process(request)
		}
	}
}
