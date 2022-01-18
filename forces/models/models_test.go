package models

import (
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
