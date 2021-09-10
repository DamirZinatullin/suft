package suft_api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAuthToken(t *testing.T) {
	suftAPI := NewSuftAPI()
	err := suftAPI.GetAuthTokens()
	fmt.Println(suftAPI.GetAccessToken())
	fmt.Println(suftAPI.GetRefreshToken())
	require.NoError(t, err)
}
