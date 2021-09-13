package http_client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthTokens(t *testing.T) {
	httpClient := NewHttpClient()
	err := httpClient.AuthTokens()
	fmt.Println("Access-token:", httpClient.GetAccessToken())
	fmt.Println("Refresh-token:", httpClient.GetRefreshToken())
	require.NoError(t, err)
}

func TestSchedules(t *testing.T) {
	httpClient := NewHttpClient()
	err := httpClient.AuthTokens()
	require.NoError(t, err)
	schedules, err := httpClient.Schedules()
	require.NoError(t, err)
	fmt.Println("schedules:", schedules)
}
