package quadratic

import (
	"fmt"
	"ga/forces/models"
	"testing"
)

func TestLinearInterpolation(t *testing.T) {
	t.Run("testing linear interpolation recursion with disparate parents", func(t *testing.T) {
		iter := 0.0

		lin, err := LinearInterpolation(&iter, 0.67, 0.31, 0.62)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(lin)

	})

	t.Run("testing too large beta", func(t *testing.T) {

	})
}
func TestGaussianFit(t *testing.T) {
	t.Run("testing quadratic fitting with random parents", func(t *testing.T) {
		a := models.Organism{
			DNA: []models.Chromosome{
				{0.648147715662861, 0.9943472031080574, 0.3180209291903924},
				{1.3574156886397406, 0.34752483290999153, -1.7787838582081046},
				{0.9080929740108076, -0.11394537066514232, -0.39033312515899127},
			},

			Path:    "",
			Fitness: 5,
		}

		b := models.Organism{
			DNA: []models.Chromosome{
				{1.0, 1.5, 2.0},
				{0.1, 0.2, 0.3},
				{-0.5, -1.0, -1.5},
			},

			Path:    "",
			Fitness: 20,
		}

		c := models.Organism{
			DNA: []models.Chromosome{
				{0.25, 0.5, 0.75},
				{1.0, 1.25, 1.50},
				{-0.65, -0.95, -1.95},
			},

			Path:    "",
			Fitness: 15,
		}

		child := QuadraticTerms(&a, &b, &c)

		fmt.Printf("%#v", child)
	})
}
