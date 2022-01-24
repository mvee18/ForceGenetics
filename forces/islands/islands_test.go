package islands

import (
	"ga/forces/models"
	"testing"
)

var (
	immPGA = make(chan models.OrganismAndBias)
	immTGA = make(chan models.OrganismAndBias)
	immIGA = make(chan models.OrganismAndBias)
)

func TestRunIslands(t *testing.T) {
	t.Run("receive from islands?", func(t *testing.T) {
		RunIslands(immPGA, immTGA, immIGA)
	})
}
