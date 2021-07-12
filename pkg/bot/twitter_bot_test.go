package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareTweetContent(t *testing.T) {
	assert.Equal(t, "The great...", prepareTweetContent("The great revolution", 9))
}
