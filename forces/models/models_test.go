package models

import (
	"math"
	"testing"
)

/*
func TestSummarize(t *testing.T) {
	f, err := os.Open("spectro.out")
	if err != nil {
		t.Fatalf("failed to open file")
	}

	result := summarize.Spectro(f)

	fmt.Printf("%#v", result)
}
*/

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

func TestCalculateBias(t *testing.T) {

	t.Run("testing with single chromosome population", func(t *testing.T) {
		var d = DNA{
			Chromosome{1.0},
			// Chromosome{4.0, 5.0},
			// Chromosome{7.0, 8.0},
		}

		var d2 = DNA{
			Chromosome{1.0},
			// Chromosome{4.0, -5.0},
			// Chromosome{-7.0, 8.0},
		}

		var dopp = DNA{
			Chromosome{-1.0},
			// Chromosome{4.0, -5.0},
			// Chromosome{-7.0, 8.0},
		}

		var o = Organism{
			DNA:     d,
			Path:    "",
			Fitness: 0.0,
		}

		var o2 = Organism{
			DNA:     d2,
			Path:    "",
			Fitness: 0.0,
		}

		var oopp = Organism{
			DNA:     dopp,
			Path:    "",
			Fitness: 0.0,
		}

		t.Run("test with exactly the same", func(t *testing.T) {
			pop := []Organism{o, o2}

			bias := CalculateBias(pop)

			want := 1.0

			if bias != want {
				t.Errorf("wrong bias, wanted %v, got %v\n", want, bias)
			}
		})

		t.Run("test with exact opposites", func(t *testing.T) {
			pop := []Organism{o, oopp}

			bias := CalculateBias(pop)

			want := 0.5

			if bias != want {
				t.Errorf("wrong bias, wanted %v, got %v\n", want, bias)
			}
		})
	})

	t.Run("testing with diff length multi-chromosomes", func(t *testing.T) {
		var d = DNA{
			Chromosome{1.0},
			Chromosome{3.0, 3.0},
			// Chromosome{7.0, 8.0},
		}

		var d2 = DNA{
			Chromosome{1.0},
			Chromosome{3.0, 3.0},
			// Chromosome{-7.0, 8.0},
		}

		var dopp = DNA{
			Chromosome{-1.0},
			Chromosome{-3.0, -3.0},
			// Chromosome{-7.0, 8.0},
		}

		var o = Organism{
			DNA:     d,
			Path:    "",
			Fitness: 0.0,
		}

		var o2 = Organism{
			DNA:     d2,
			Path:    "",
			Fitness: 0.0,
		}

		var oopp = Organism{
			DNA:     dopp,
			Path:    "",
			Fitness: 0.0,
		}

		t.Run("testing exactly the same", func(t *testing.T) {
			pop := []Organism{o, o2}

			bias := CalculateBias(pop)

			want := 1.0

			if bias != want {
				t.Errorf("wrong bias, wanted %v, got %v\n", want, bias)
			}
		})

		t.Run("testing exactly the opposite", func(t *testing.T) {
			pop := []Organism{o, oopp}

			bias := CalculateBias(pop)

			want := 0.5

			if bias != want {
				t.Errorf("wrong bias, wanted %v, got %v\n", want, bias)
			}
		})
	})

}

func TestCalculateHD(t *testing.T) {
	var d = DNA{
		Chromosome{1.0},
		Chromosome{3.0, 3.0},
		// Chromosome{7.0, 8.0},
	}

	var d2 = DNA{
		Chromosome{1.0},
		Chromosome{3.0, 3.0},
		// Chromosome{-7.0, 8.0},
	}

	var dopp = DNA{
		Chromosome{-1.0},
		Chromosome{-3.0, -3.0},
		// Chromosome{-7.0, 8.0},
	}

	var o = Organism{
		DNA:     d,
		Path:    "",
		Fitness: 0.0,
	}

	var o2 = Organism{
		DNA:     d2,
		Path:    "",
		Fitness: 0.0,
	}

	var oopp = Organism{
		DNA:     dopp,
		Path:    "",
		Fitness: 0.0,
	}

	t.Run("testing HD when the same", func(t *testing.T) {
		hd := CalculateHD(o, o2)

		wanthd := 0.0

		if hd != wanthd {
			t.Errorf("wrong hamming distance, wanted %v, got %v\n", wanthd, hd)
		}

	})

	t.Run("testing HD when the exact opposite (maximum)", func(t *testing.T) {
		hd := CalculateHD(o, oopp)

		wanthd := 3.0

		if hd != wanthd {
			t.Errorf("wrong hamming distance, wanted %v, got %v\n", wanthd, hd)
		}

	})
}

func TestPopulationAttractiveness(t *testing.T) {
	o1 := Organism{
		DNA:     DNA{},
		Fitness: 1.0,
		Path:    "",
	}

	t.Run("check if pop attractives calc correctly", func(t *testing.T) {
		newFitness := []float64{1.0, 6.0, 10.0}

		mig := Migrant{
			Attractiveness: 0.0,
			Org:            o1,
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
		o1 := Organism{
			DNA:     DNA{},
			Fitness: 1.0,
			Path:    "",
		}

		mig := Migrant{
			Attractiveness: 0.0,
			Org:            o1,
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
		o1 := Organism{
			DNA:     DNA{},
			Fitness: 1.0,
			Path:    "",
		}

		newFitness := []float64{1.0, 6.0, 10.0}

		// initialize with 5 fitness, should be added to previous.
		mig := Migrant{
			Attractiveness: 5.0,
			Org:            o1,
			PrevFitness:    []float64{2.0, 12.0, 20.0},
		}

		mig.CalculateAttractiveness(newFitness)

		want := (5.0 + 17.0/3.0 + 1.0)

		// Some floating point precision loss makes them not exactly equal.
		// Just use the difference to 8 decimals. That's more than close enough.
		if math.Abs(mig.Attractiveness-want) > 1e-8 {
			t.Errorf("error calc attractivness, got %v, wanted %v\n", mig.Attractiveness, want)
		}

	})
}
