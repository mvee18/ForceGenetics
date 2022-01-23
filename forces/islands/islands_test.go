package islands

import (
	"ga/forces/models"
	"testing"
)

func TestRunIslands(t *testing.T) {
	t.Run("receive from islands?", func(t *testing.T) {
		imm, mig := make(chan models.Organism), make(chan models.Organism)

		RunIslands(imm, mig)

	})
}
