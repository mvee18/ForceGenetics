package pseudo

import (
	"testing"
)

func TestGeneratePseudoPopulation(t *testing.T) {
	t.Run("testing if pop generated w/o overflow", func(t *testing.T) {
		CreatePseudoPopulation()
	})
}
