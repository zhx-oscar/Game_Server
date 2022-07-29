package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureIndexes(t *testing.T) {
	err := EnsureIndexes()
	assert.Nil(t, err)
}
