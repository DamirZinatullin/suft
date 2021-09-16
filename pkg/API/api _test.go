package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchedules(t *testing.T) {
	apiClient, err := NewClient("demo@example.com", "demo")
	require.NoError(t, err)
	options := &Options{
		Page:            1,
		Size:            5,
		CreatorApprover: "creator",
	}
	schedules, err := apiClient.Schedules(options)
	require.NoError(t, err)
	for _, schedule := range schedules {
		fmt.Println(schedule)
	}
}
