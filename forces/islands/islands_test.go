package islands

import (
	"ga/forces/models"
	"math"
	"testing"
)

// var (
// 	immPGA = make(chan models.OrganismAndBias)
// 	immTGA = make(chan models.OrganismAndBias)
// 	immIGA = make(chan models.OrganismAndBias)
// )

// func TestRunIslands(t *testing.T) {
// 	t.Run("receive from islands?", func(t *testing.T) {
// 		RunIslands(immPGA, immTGA, immIGA)
// 	})
// }

func TestPopulationAttractiveness(t *testing.T) {
	o1 := models.Organism{
		DNA:     models.DNA{},
		Fitness: 1.0,
		Path:    "",
	}

	t.Run("check if pop attractives calc correctly", func(t *testing.T) {
		newFitness := []float64{1.0, 6.0, 10.0}

		mig := models.Migrant{
			Attractiveness: 0.0,
			Mig:            o1,
			PrevFitness:    []float64{2.0, 12.0, 20.0},
		}

		got := populationAttractiveness(&mig, newFitness)

		want := 17.0 / 3.0

		if got != want {
			t.Errorf("wrong pop fitness, wanted %v got %v\n", want, got)
		}
	})
}

func TestMigrantAttractiveness(t *testing.T) {
	t.Run("testing calc of mig attractiveness", func(t *testing.T) {
		o1 := models.Organism{
			DNA:     models.DNA{},
			Fitness: 1.0,
			Path:    "",
		}

		mig := models.Migrant{
			Attractiveness: 0.0,
			Mig:            o1,
			PrevFitness:    []float64{2.0, 12.0, 20.0},
		}

		got := migrantAttractiveness(&mig)

		want := 1.0

		if got != want {
			t.Errorf("error calc migA, got %v wanted %v\n", got, want)
		}
	})
}

func TestCalcAttractiveness(t *testing.T) {
	t.Run("test calc overall attractivness", func(t *testing.T) {
		o1 := models.Organism{
			DNA:     models.DNA{},
			Fitness: 1.0,
			Path:    "",
		}

		newFitness := []float64{1.0, 6.0, 10.0}

		// initialize with 5 fitness, should be added to previous.
		mig := models.Migrant{
			Attractiveness: 5.0,
			Mig:            o1,
			PrevFitness:    []float64{2.0, 12.0, 20.0},
		}

		CalculateAttractiveness(&mig, newFitness)

		want := (5.0 + 17.0/3.0 + 1.0)

		// Some floating point precision loss makes them not exactly equal.
		// Just use the difference to 8 decimals. That's more than close enough.
		if math.Abs(mig.Attractiveness-want) > 1e-8 {
			t.Errorf("error calc attractivness, got %v, wanted %v\n", mig.Attractiveness, want)
		}

	})
}
