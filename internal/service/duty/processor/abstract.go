package processor

import "svtt/internal/entities"

type Dutyer interface {
	Process(request entities.DutyRequest)
	Stop()
}

type SingleDutyProcessor func(request entities.DutyRequest) entities.DutyResponce
