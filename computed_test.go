package reactor_test

import (
	"reactor"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputed(t *testing.T) {
	t.Run("computed is updated", func(t *testing.T) {
		a := reactor.NewSignal(2)
		require.Equal(t, 2, a.Get())

		b := reactor.NewSignal(6)
		require.Equal(t, 6, b.Get())

		c := reactor.NewComputed(func() int {
			return a.Get() + b.Get()
		})
		require.Equal(t, 2+6, c.Get())

		b.Set(12)
		a.Set(22)
		require.Equal(t, 12+22, c.Get())

		b.Set(12)
		a.Set(0)
		require.Equal(t, 12, c.Get())
	})
}
