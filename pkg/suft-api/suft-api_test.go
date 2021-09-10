package suft_api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAuthToken(t *testing.T) {
	suftAPI := NewSuftAPI()
	err := suftAPI.GetAuthToken()
	fmt.Println(suftAPI.AccessToken)
	fmt.Println(suftAPI.RefreshToken)
	require.NoError(t, err)
}
