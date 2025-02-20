package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword(t *testing.T) {
	password := "abcdefg"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, wrongPassword)
	require.NotEqual(t, wrongPassword, hashedPassword)
}
