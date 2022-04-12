package quadratic

import (
	"errors"
	"ga/forces/models"
	"ga/forces/utils"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Suppose that we have three parents. We can fit a quadratic equation using the following terms from each:

var r1 *rand.Rand

func init() {
	rand.Seed(time.Now().UnixNano())
	s1 := rand.NewSource(time.Now().UnixNano() + 2561)
	r1 = rand.New(s1)
}

var ErrLinearFailed = errors.New("maximum number of iterations reached with reduction of beta")

func QuadraticTerms(o1, o2, o3 *models.Organism) models.Organism {

	p1, p2, p3 := reOrderOrganisms([]*models.Organism{o1, o2, o3})

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

			maximum, valid := calcMaximum(i+2, aj, bj)
			if !valid {
				beta := r1.Float64()
				linear, err := LinearInterpolation(utils.SelectDomain(i+2), beta, p1.DNA[i][j], p3.DNA[i][j])
				if err == ErrLinearFailed {
					// fmt.Println(err)
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

// The fitness function is actually 1/F, since a lower fitness is better. By
// finding the maximum fitness, we are actively working against the process.

func reOrderOrganisms(o []*models.Organism) (models.Organism, models.Organism, models.Organism) {
	p := make([]models.Organism, len(o))

	for i := range p {
		p[i].DNA = (o)[i].DNA
		// p[i].Fitness = 1 / (o)[i].Fitness
		p[i].Fitness = (o)[i].Fitness
	}

	sort.SliceStable(p, func(i, j int) bool {
		return p[i].Fitness < p[j].Fitness
	})

	return p[0], p[1], p[2]

}

func calcMaximum(dn int, aj, bj float64) (float64, bool) {
	Ej := -bj / (2 * aj)
	// fmt.Printf("Ej is %v, and aj is %v\n", Ej, aj)

	if 2*aj > 0 && math.Abs(Ej) < utils.SelectDomain(dn) {
		return Ej, true
	} else {
		return Ej, false
	}
}

// These need to be sorted.
func LinearInterpolation(alpha, beta, m, d float64) (float64, error) {
	// CrossoverPoint
	// for *iterations <= 3 {
	// 	pNew := beta*(m-d) + m

	// 	if math.Abs(pNew) < alpha {
	// 		*iterations = 0.0
	// 		return pNew, nil

	// 	} else {
	// 		bNew := beta / 2
	// 		lin, err := LinearInterpolation(iterations, alpha, bNew, m, d)
	// 		*iterations = *iterations + 1.0
	// 		return lin, err
	// 	}
	// }
	// 	pNew := beta*(m-d) + m
	iterations := new(int)

	pNew := beta*(m-d) + m
	// fmt.Printf("the pnew is %v\n", pNew)

	for math.Abs(pNew) > alpha {
		bNew := beta / 2
		pNew = bNew*(m-d) + m

		*iterations = *iterations + 1
		// fmt.Printf("the number of iterations is %v\n", *iterations)
		if *iterations > 3 {
			return 0.0, ErrLinearFailed
		}
	}

	return pNew, nil
}
