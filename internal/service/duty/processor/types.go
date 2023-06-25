package processor

import "svtt/internal/entities"

const (
	processingStatusReady = iota + 1
	processingStatusRunning
	processingStatusDone
)

type processingContainer struct {
	status  int
	request entities.DutyRequest
}
