package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/ntBre/chemutils/summarize"
)

func TestSummarize(t *testing.T) {
	f, err := os.Open("spectro.out")
	if err != nil {
		t.Fatalf("failed to open file")
	}

	result := summarize.Spectro(f)

	fmt.Printf("%#v", result)
}

func BenchmarkCalcFitness(t *testing.B) {
	t.Run("how slow is calc fitness?", func(b *testing.B) {
		d := Organism{
			DNA:     DNA{},
			Fitness: 0.0,
			Path:    "../informed/testorg/fort.40",
		}

		for n := 0; n < b.N; n++ {
			d.CalcFitness()
		}

	})
}
