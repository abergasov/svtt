package duty_processor

import "svtt/internal/logger"

type Service struct {
	log logger.AppLogger
}

func NewService(log logger.AppLogger) *Service {
	return &Service{log: log}
}
