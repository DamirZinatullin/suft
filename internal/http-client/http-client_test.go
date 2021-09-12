package http_client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAuthToken(t *testing.T) {
	httpClient := NewHttpClient()
	err := httpClient.AuthTokens()
	fmt.Println(httpClient.GetAccessToken())
	fmt.Println(httpClient.GetRefreshToken())
	require.NoError(t, err)
}
