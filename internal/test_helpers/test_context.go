package testhelpers

import (
	"svtt/internal/config"
	"svtt/internal/logger"
	"svtt/internal/service/duty_processor"

	"testing"

	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	ServiceDutyProcessor *duty_processor.Service
}

func GetClean(t *testing.T) *TestContainer {
	// conf := getTestConfig()

	appLog, err := logger.NewAppLogger("test")
	require.NoError(t, err)

	// service init
	serviceDuty := duty_processor.NewService(appLog)
	return &TestContainer{
		ServiceDutyProcessor: serviceDuty,
	}
}

func getTestConfig() *config.AppConfig {
	return &config.AppConfig{
		AppPort: 0,
	}
}
