package reactor_test

import (
	"reactor"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEffect(t *testing.T) {
	t.Run("effect is triggered when dependent signal is updated", func(t *testing.T) {
		signal := reactor.NewSignal(6)

		computed := reactor.NewComputed(func() int {
			return signal.Get() * 2
		})
		require.Equal(t, 12, computed.Get())

		effectExecutionCounter := 0
		reactor.NewEffect(func() {
			_ = signal.Get()
			effectExecutionCounter++
		})
		require.Equal(t, 1, effectExecutionCounter)

		signal.Set(9)
		require.Equal(t, 2, effectExecutionCounter)
		require.Equal(t, 18, computed.Get())

		signal.Set(12)
		require.Equal(t, 3, effectExecutionCounter)
		require.Equal(t, 24, computed.Get())
	})
}
