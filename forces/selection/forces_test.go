package selection

import (
	"fmt"
	"ga/forces/models"
	"ga/forces/utils"
	"math"
	"testing"
)

func TestCreateOrganism(t *testing.T) {
	t.Run("run org generation", func(t *testing.T) {
		organism := CreateOrganism(3)
		if len(organism.DNA[0]) != utils.GetNumForceConstants(3, 2) {
			t.Errorf("organism was not the right size")
		}

		if organism.Fitness == 0 {
			t.Errorf("Error calculating fitness\n")
		}

		fmt.Println(organism.Fitness)
	})
}

func TestCalcFitness(t *testing.T) {
	t.Run("checking if LXM is correct for known fort.15 h2o", func(t *testing.T) {
		organism := models.Organism{
			DNA:     nil,
			Path:    "testfiles/h2o/2nd/fort.15",
			Fitness: 0,
		}

		organism.CalcFitness()

		want := 3.218
		got := organism.Fitness

		if got >= want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}
	})

	t.Run("running on bad organism LXM", func(t *testing.T) {
		org := models.Organism{
			DNA:     nil,
			Path:    "testfiles/bad/20/fort.15",
			Fitness: 0,
		}

		org.CalcFitness()

		want := fmt.Sprintf("%.4f", math.Sqrt(17426396.1))
		got := fmt.Sprintf("%.4f", org.Fitness)

		if got != want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}

	})

	/*
		t.Run("check 3rd rotational constants", func(t *testing.T) {
			organism := Organism{
				DNA:     nil,
				Path:    "testfiles/h2o/3rd/fort.30",
				Fitness: 0,
			}

			organism.calcFitness(TargetFrequencies)

			want := 0.0
			got := organism.Fitness

			if got != want {
				t.Errorf("got %v, wanted %v\n", want, got)
			}

		})
	*/
}

func TestCreatePopulation(t *testing.T) {
	t.Run("run create pop", func(t *testing.T) {
		pop := CreatePopulation()
		pool := CreatePool(pop, models.TargetFrequencies)

		population := NaturalSelection(pool, pop, models.TargetFrequencies)

		fmt.Println(len(population))
	})
}

func TestGetNumForceConstants(t *testing.T) {
	t.Run("testing number of force constants for water", func(t *testing.T) {
		got := utils.GetNumForceConstants(3, 3)
		want := 165

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}

	})

	t.Run("testing 4th derivative water", func(t *testing.T) {
		got := utils.GetNumForceConstants(3, 4)
		want := 495

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("testing number for 6 atoms", func(t *testing.T) {
		got := utils.GetNumForceConstants(6, 3)
		want := 1140

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}

func TestReadInput(t *testing.T) {
	t.Run("testing testfile input", func(t *testing.T) {
		utils.ReadInput("testfiles/forces.inp")
	})
}
