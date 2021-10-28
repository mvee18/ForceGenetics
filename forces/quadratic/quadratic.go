package quadratic

import (
	"errors"
	"ga/forces/models"
	"math"
	"math/rand"
)

// Suppose that we have three parents. We can fit a quadratic equation using the following terms from each:

const alpha = 2

var ErrLinearFailed = errors.New("maximum number of iterations reached with reduction of beta")

func QuadraticTerms(p1, p2, p3 *models.Organism) models.Organism {
	childDNA := make(models.DNA, len(p1.DNA))

	for i, pChr := range p1.DNA {
		ch := make(models.Chromosome, len(pChr))
		childDNA[i] = ch
	}

	child := models.Organism{
		DNA:     childDNA,
		Path:    "",
		Fitness: 0,
	}

	for i, chromosome := range p1.DNA {
		for j := range chromosome {
			aj := 1 / (p3.DNA[i][j] - p2.DNA[i][j]) * (((p3.Fitness - p1.Fitness) / (p3.DNA[i][j] - p1.DNA[i][j])) - ((p2.Fitness - p1.Fitness) / (p2.DNA[i][j] - p1.DNA[i][j])))

			bj := ((p2.Fitness - p1.Fitness) / (p2.DNA[i][j] - p1.DNA[i][j])) - aj*(p2.DNA[i][j]+p1.DNA[i][j])

			// cj := p1.Fitness - aj*math.Pow(p1.DNA[i][j], 2) - bj*p1.DNA[i][j]

			// fmt.Println(aj, bj, cj)

			maximum, valid := calcMaximum(aj, bj)
			if !valid {
				beta := rand.Float64()
				iterations := 0.0
				linear, err := LinearInterpolation(&iterations, beta, p1.DNA[i][j], p3.DNA[i][j])
				if err == ErrLinearFailed {
					p := rand.Intn(3)
					switch p {
					case 0:
						child.DNA[i][j] = p1.DNA[i][j]

					case 1:
						child.DNA[i][j] = p2.DNA[i][j]

					case 2:
						child.DNA[i][j] = p3.DNA[i][j]

					}

				} else {
					child.DNA[i][j] = linear
				}

			} else {
				child.DNA[i][j] = maximum
			}
		}
	}

	return child
}

func calcMaximum(aj, bj float64) (float64, bool) {
	Ej := -bj / (2 * aj)

	if 2*aj < 0 && math.Abs(Ej) < alpha {
		return Ej, true
	} else {
		return Ej, false
	}
}

// These need to be sorted.
func LinearInterpolation(iterations *float64, beta, m, d float64) (float64, error) {
	// CrossoverPoint
	if *iterations > 3 {
		return 0.0, ErrLinearFailed
	}

	pNew := beta*(m-d) + m

	if math.Abs(pNew) < alpha {
		return pNew, nil

	} else {
		bNew := beta / 2
		lin, err := LinearInterpolation(iterations, bNew, m, d)
		*iterations = *iterations + 1.0
		return lin, err
	}

}
