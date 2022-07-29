package dbidx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureIndexes(t *testing.T) {
	err := EnsureIndexes()
	require.NoError(t, err)
}
