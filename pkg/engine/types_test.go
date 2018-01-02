package engine

import "testing"
import "github.com/stretchr/testify/assert"

func TestTypes(t *testing.T) {
	t.Parallel()
	t.Run("square-towards", func(tt *testing.T) {
		assert.Equal(tt, B1, C1.Towards(West))
	})
}
