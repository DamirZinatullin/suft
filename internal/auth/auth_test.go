package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	token, err := Authenticate("demo@example.com", "demo")
	require.NoError(t, err)
	fmt.Println("Access-token:", token.AccessToken)
	fmt.Println("Refresh-token:", token.RefreshToken)
}
