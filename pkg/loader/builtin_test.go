package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltins(t *testing.T) {
	assert.NotNil(t, Builtins)
	assert.Greater(t, len(Builtins), 0)
}
