package builtins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbbedChecks(t *testing.T) {
	entries, err := EmbbedChecksFS.ReadDir(".")
	assert.NoError(t, err)
	assert.Greater(t, len(entries), 0)
}
