package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicKey_IsCaseInsensitive(t *testing.T) {
	a, err := NewPublicKey("ABCD")
	require.NoError(t, err)

	b, err := NewPublicKey("abcd")
	require.NoError(t, err)

	require.Equal(t, a, b)
}
