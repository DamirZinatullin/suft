package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAuthenticateSuccess(t *testing.T) {
	token, err := Authenticate("demo@example.com", "demo", nil)
	require.NoError(t, err)
	fmt.Println("Access-token:", token.AccessToken)
	fmt.Println("Refresh-token:", token.RefreshToken)
}

func TestAuthenticateUnauthorized(t *testing.T) {
	token, err := Authenticate("fake@example.com", "fake", nil)
	require.Error(t, err)
	assert.Nil(t, token)
}

func TestRefresh(t *testing.T) {
	token, err := Authenticate("demo@example.com", "demo", nil)
	require.NoError(t, err)
	tokenResp, err := Refresh(token.RefreshToken, nil)
	require.NoError(t, err)
	fmt.Println("Access-token:", tokenResp.AccessToken)
	fmt.Println("Refresh-token:", tokenResp.RefreshToken)
}

func TestRefreshFail(t *testing.T) {
	token, err := Refresh("fake", nil)
	require.Error(t, err)
	assert.Nil(t, token)
}
