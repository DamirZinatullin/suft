package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchedules(t *testing.T) {
	apiClient, err := NewClient("demo@example.com", "demo")
	require.NoError(t, err)
	schedules, err := apiClient.Schedules(nil)
	require.NoError(t, err)
	for _, schedule := range schedules {
		fmt.Println(schedule)
	}
}

func TestLoggingTimeList(t *testing.T) {
	apiClient, err := NewClient("demo@example.com", "demo")
	require.NoError(t, err)
	loggingTimes, err := apiClient.LoggingTimeList(32907, nil)
	require.NoError(t, err)
	for _, loggingTime := range loggingTimes {
		fmt.Println(loggingTime)
	}
}
