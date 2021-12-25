package pseudo

import (
	"ga/forces/models"
	"reflect"
	"testing"
)

var d = models.DNA{
	models.Chromosome{2.57, 52.25},
	models.Chromosome{-68.35, 1.23},
	models.Chromosome{-1.1, 1.2, 1.3},
}

var dcomp = models.DNA{
	models.Chromosome{-2.57, -52.25},
	models.Chromosome{68.35, -1.23},
	models.Chromosome{1.1, -1.2, -1.3},
}

var d2 = models.DNA{
	models.Chromosome{-1.0, 2.0, 3.0},
	models.Chromosome{4.0, -5.0, 6.0},
	models.Chromosome{7.0, 8.0, -9.0},
}

var d2comp = models.DNA{
	models.Chromosome{1.0, -2.0, -3.0},
	models.Chromosome{-4.0, 5.0, -6.0},
	models.Chromosome{-7.0, -8.0, 9.0},
}

var o = models.Organism{
	DNA:     d,
	Path:    "",
	Fitness: 0.0,
}

var ocomp = models.Organism{
	DNA:     dcomp,
	Path:    "",
	Fitness: 0.0,
}

var o2 = models.Organism{
	DNA:     d2,
	Path:    "",
	Fitness: 0.0,
}

var o2comp = models.Organism{
	DNA:     d2comp,
	Path:    "",
	Fitness: 0.0,
}

func TestGeneratePseudoPopulation(t *testing.T) {
	t.Run("testing if pop generated w/o overflow", func(t *testing.T) {
		pop := CreatePseudoPopulation()

		t.Run("verify if adjacent organism are complimentary", func(t *testing.T) {
			for i := 0; i < len(pop); i += 2 {
				org := pop[i]

				orgComp := complimentaryOrganism(&org)

				if !reflect.DeepEqual(orgComp.DNA, pop[i+1].DNA) {
					t.Errorf("adjacent DNA not complimentary in organism index %d\n", i)
				}

			}

		})

	})
}

func TestComplimentaryOrganism(t *testing.T) {
	t.Run("generate complimentary organism", func(t *testing.T) {
		comp := complimentaryOrganism(&o)

		if !reflect.DeepEqual(comp.DNA, ocomp.DNA) {
			t.Errorf("DNA not complimentary")
		}
	})
}

func TestPseudoCrossover(t *testing.T) {
	var d = models.DNA{
		models.Chromosome{2.57, 52.25},
		models.Chromosome{-68.35, 1.23},
		models.Chromosome{-1.1, 1.2, 1.3},
	}

	var dcomp = models.DNA{
		models.Chromosome{-2.57, -52.25},
		models.Chromosome{68.35, -1.23},
		models.Chromosome{1.1, -1.2, -1.3},
	}

	var o = models.Organism{
		DNA:     d,
		Path:    "",
		Fitness: 0.0,
	}

	var ocomp = models.Organism{
		DNA:     dcomp,
		Path:    "",
		Fitness: 0.0,
	}

	t.Run("test perform crossover", func(t *testing.T) {
		// Crossover at 100% btwn two compliments should yield the
		// opposite.

		c1, c2 := performPseudoCrossover(o, ocomp)

		if !reflect.DeepEqual(c1.DNA, ocomp.DNA) || !reflect.DeepEqual(c2.DNA, o.DNA) {
			t.Errorf("DNA not complimentary in organisms %v %v\n", c1.DNA, ocomp.DNA)
		}

	})

	t.Run("test if correct organisms made with 100 crossrate", func(t *testing.T) {
		var pop = PPopulation{o, ocomp, o2, o2comp}

		newPop := psuedoCrossover(pop)

		// If this crossover happened correctly at 100%, then the two
		// organisms should have swapped places. Thus, the original
		// population's first element should be equal to the second
		// element in the new population.
		for i := 0; i < len(newPop); i += 2 {
			if !reflect.DeepEqual(newPop[i].DNA, pop[i+1].DNA) {
				t.Errorf("error in complimentarity in organism %v and %v\n", newPop[i].DNA, pop[i-1].DNA)
			}
		}

	})
}

func TestRunPseudoGA(t *testing.T) {
	t.Run("testing PGA", func(t *testing.T) {
		runPGA()
	})
}
