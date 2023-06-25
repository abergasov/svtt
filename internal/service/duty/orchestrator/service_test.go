package orchestrator_test

import (
	"svtt/internal/entities"
	"svtt/internal/logger"
	"svtt/internal/service/duty/orchestrator"
	"svtt/internal/service/duty/processor"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	sleepSeconds = 2
	heightsCount = 4 // how many heights will be processed per validator and per duty type
)

func TestService_HandleDutyRequest(t *testing.T) {
	// given
	appLog, err := logger.NewAppLogger("test")
	require.NoError(t, err)

	// generate validator duty requests
	knowledge := make(map[int]map[entities.DutyType]map[int]entities.DutyResponce)
	executionOrder := make(map[int]map[entities.DutyType]int)      // track execution order
	currentRunning := make(map[int]map[entities.DutyType]struct{}) // track currentrly running duties
	stateMU := sync.Mutex{}
	dutyes := []entities.DutyType{
		entities.DutyTypeProposer,
		entities.DutyTypeAggregator,
		entities.DutyTypeSyncCommittee,
		entities.DutyTypeAttester,
	}
	counter := int32(0)
	for i := 0; i < 10; i++ { // validators
		knowledge[i] = make(map[entities.DutyType]map[int]entities.DutyResponce)
		currentRunning[i] = make(map[entities.DutyType]struct{})
		for j := 0; j < len(dutyes); j++ { // duty type
			knowledge[i][dutyes[j]] = make(map[int]entities.DutyResponce)
			for k := 1; k <= heightsCount; k++ { // height
				knowledge[i][dutyes[j]][k] = entities.DutyResponce{
					Duty:      dutyes[j],
					Height:    k,
					Validator: i,
					Response:  uuid.NewString(),
				}
				counter++
			}
		}
	}
	duplicates := make(map[string]struct{})

	// when
	service := orchestrator.NewService(appLog, processor.NewService, orchestrator.WithCustomRequestProcessor(func(request entities.DutyRequest) entities.DutyResponce {
		// check that only one processor per dutyTipe executes at the same time
		stateMU.Lock()
		if _, ok := currentRunning[request.Validator][request.Duty]; ok {
			t.Fatal("duty already running")
		}
		currentRunning[request.Validator][request.Duty] = struct{}{}
		previousHeigt := executionOrder[request.Validator][request.Duty]
		if previousHeigt >= request.Height {
			t.Fatal("wrong execution order")
		}
		stateMU.Unlock()

		defer func() {
			stateMU.Lock()
			delete(currentRunning[request.Validator], request.Duty)
			stateMU.Unlock()
		}()

		time.Sleep(sleepSeconds * time.Second)
		return entities.DutyResponce{
			Duty:      request.Duty,
			Height:    request.Height,
			Validator: request.Validator,
			Response:  knowledge[request.Validator][request.Duty][request.Height].Response,
		}
	}))
	t.Cleanup(func() {
		require.NoError(t, service.Stop())
	})

	// listen for duty responses and compare with knowledge
	go func() {
		for res := range service.GetDutyResponser() {
			atomic.AddInt32(&counter, -1)
			validator, ok := knowledge[res.Validator]
			require.True(t, ok)
			duty, ok := validator[res.Duty]
			require.True(t, ok)
			req, ok := duty[res.Height]
			require.True(t, ok)
			require.Equal(t, req.Response, res.Response)
			_, ok = duplicates[res.Response]
			require.False(t, ok)
			duplicates[res.Response] = struct{}{}
		}
	}()

	// send duty requests
	for validator, dutiesTypes := range knowledge {
		for duty, duties := range dutiesTypes {
			for height := range duties {
				service.HandleDutyRequest(entities.DutyRequest{
					Duty:      duty,
					Height:    height,
					Validator: validator,
				})
			}
		}
	}

	// then
	waitTime := time.Duration((heightsCount * sleepSeconds * len(dutyes)) + 1) // 1 - just in case for goroutines switching
	require.Eventuallyf(t, func() bool {
		return atomic.LoadInt32(&counter) == 0
	}, waitTime*time.Second, 100*time.Millisecond, "not all duty requests were processed: %d", atomic.LoadInt32(&counter))
}
