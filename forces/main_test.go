package main

import (
	"fmt"
	"ga/forces/models"
	"ga/forces/selection"
	"testing"
)

func TestGenerateOrganism(t *testing.T) {
	t.Run("testing generate organism", func(t *testing.T) {
		org := selection.CreateOrganism(4)

		fmt.Println(org.Fitness)
	})
}

func TestCalcFitness(t *testing.T) {
	t.Run("Testing formaldehyde with qff force constants", func(t *testing.T) {
		known := models.Organism{
			DNA:     nil,
			Fitness: 0,
			Path:    "testfiles/h2co/4th/",
		}

		known.CalcFitness()

		if known.Fitness > 0.01 {
			t.Errorf("wanted 0, got %v\n", known.Fitness)
		}
	})
}
