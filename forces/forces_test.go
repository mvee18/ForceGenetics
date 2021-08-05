package main

import (
	"fmt"
	"testing"
)

func TestCreateOrganism(t *testing.T) {
	t.Run("run org generation", func(t *testing.T) {
		organism := CreateOrganism(6)
		if len(organism.DNA) != GetNumForceConstants(6, 2) {
			t.Errorf("organism was not the right size")
		}

		if organism.Fitness == 0 {
			t.Errorf("Error calculating fitness\n")
		}

		fmt.Println(organism.Fitness)
	})
}

func TestCalcFitness(t *testing.T) {
	t.Run("checking if LXM is correct for known fort.15", func(t *testing.T) {
		organism := Organism{
			DNA:     nil,
			Path:    "testfiles/fort.15",
			Fitness: 0,
		}

		organism.calcFitness(TargetFrequencies)

		want := 0.0
		got := organism.Fitness

		if got != want {
			t.Errorf("got %v, wanted %v\n", want, got)
		}
	})
}

func TestCreatePopulation(t *testing.T) {
	t.Run("run create pop", func(t *testing.T) {
		pop := createPopulation()
		pool := createPool(pop, TargetFrequencies)
		population := naturalSelection(pool, pop, TargetFrequencies)

		fmt.Println(population)
	})
}
