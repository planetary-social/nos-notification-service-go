package domain_test

import (
	"strings"
	"testing"

	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/stretchr/testify/require"
)

func TestPublicKey_IsCaseInsensitive(t *testing.T) {
	hex := fixtures.SomeHexBytesOfLen(32)
	hexLower := strings.ToLower(hex)
	hexUpper := strings.ToUpper(hex)

	require.NotEqual(t, hexLower, hexUpper)

	a, err := domain.NewPublicKey(hexLower)
	require.NoError(t, err)

	b, err := domain.NewPublicKey(hexUpper)
	require.NoError(t, err)

	require.Equal(t, a, b)
}
