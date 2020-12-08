package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeName(t *testing.T) {
	assert.Equal(t, "-test----", sanitizeName("  test - , ,"))
	assert.Equal(t, "unknown", sanitizeName(""))
}
