package processor

import (
	"svtt/internal/entities"
)

func (s *Service) processBackground() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			s.checkQueue()
			go s.cleanQueue()
		}
	}
}

// checkQueue get first avaliable duty from queue and start processing
func (s *Service) checkQueue() {
	s.executionQueueMU.Lock()
	defer s.executionQueueMU.Unlock()
	for dutyType, queue := range s.executionQueue {
		if queue.Len() == 0 {
			continue
		}
		for el := queue.Front(); el != nil; el = el.Next() {
			container := el.Value.(*processingContainer)
			if container.status == processingStatusRunning {
				break
			}
			if container.status == processingStatusReady {
				container.status = processingStatusRunning
				s.wg.Add(1)
				go s.processDuty(dutyType, container.request)
				break
			}
		}
	}
}

// processDuty process single duty request and mark it as done
func (s *Service) processDuty(dutyType entities.DutyType, request entities.DutyRequest) {
	defer s.wg.Done()
	s.responseChan <- s.processor(request)
	s.executionQueueMU.Lock()
	defer s.executionQueueMU.Unlock()
	queue := s.executionQueue[dutyType]
	for el := queue.Front(); el != nil; el = el.Next() {
		container := el.Value.(*processingContainer)
		if container.request.Height == request.Height {
			container.status = processingStatusDone
			break
		}
	}
}

// cleanQueue remove all done duties from queue
func (s *Service) cleanQueue() {
	s.executionQueueMU.Lock()
	defer s.executionQueueMU.Unlock()
	for _, queue := range s.executionQueue {
		for el := queue.Front(); el != nil; el = el.Next() {
			container := el.Value.(*processingContainer)
			if container.status == processingStatusDone {
				queue.Remove(el)
			}
		}
	}
}
