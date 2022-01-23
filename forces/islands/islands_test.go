package islands

import (
	"ga/forces/models"
	"testing"
)

var (
	immPGA = make(chan models.Organism)
	immTGA = make(chan models.Organism)
	immIGA = make(chan models.Organism)
)

func TestRunIslands(t *testing.T) {
	t.Run("receive from islands?", func(t *testing.T) {
		RunIslands(immPGA, immTGA, immIGA)
	})
}
