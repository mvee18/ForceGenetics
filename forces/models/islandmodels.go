package models

import (
	"ga/forces/utils"
	"math"
)

// Unfortunately, we have multiple chromosomes.
// As a result, we have to modify the version presented by Tsutshi, 1995.
// This will give a value 0.5 <= B <= 1.0.
func CalculateBias(p []Organism) float64 {
	// The formula specifies we compare the same gene across all the organisms
	// in the population first, not each organism.
	n := len(p)
	final := 0.0
	dna := p[0].DNA

	// Iterate over the chromosomes of the first organism. Since they're all
	// the same size, it doesn't matter which one we use.
	for j, chromosome := range dna {
		// We need all of the chromosomes in this column.
		summation := 0.0
		lChr := len(chromosome)

		for k := range chromosome {
			// Iterate again across all of the organisms (rows).
			for m := range p {
				biasValue(&summation, utils.SelectDomain(j+2), p[m].DNA[j][k], n)
				// biasValue(&summation, 3, p[m].DNA[j][k], n)
				// fmt.Printf("the domain is %v\n", utils.SelectDomain(j+2))
				// fmt.Printf("summation int: %v\n", summation)
			}

			biasSummation(&summation, n)

			// fmt.Printf("summation: %v, n: %v, lChr: %v\n", summation, n, lChr)
		}

		final += summation / float64(n) / float64(lChr)
	}

	return final / float64(len(dna))
}

func biasValue(count *float64, ub, gene float64, n int) {
	// fmt.Printf("gene: %v, ub: %v\n", gene, ub)
	*count += ((gene + ub) / (2 * ub))
}

func biasSummation(intermediate *float64, n int) {
	*intermediate -= float64(n) / 2

	*intermediate = math.Abs(*intermediate)
	*intermediate += float64(n) / 2

	// fmt.Printf("int out: %v\n", *intermediate)
}

func CalculateHD(a, b Organism) float64 {
	totalHD := 0.0

	for i, chr := range a.DNA {
		ub := utils.SelectDomain(i + 2)
		chrHD := 0.0

		for j := range chr {
			chrHD += ((a.DNA[i][j] + ub) / (2 * ub)) - ((b.DNA[i][j] + ub) / (2 * ub))
		}

		totalHD += math.Abs(chrHD)
	}

	return totalHD
}

func AddImmigrant(p *[]Organism, migrant <-chan Migrant) {
	// Take the last organism (least fit) off.
	*p = (*p)[0 : len(*p)-1]

	org := <-migrant

	*p = append(*p, org.Org)
}

func AddMigrant(migrant chan<- Migrant, best Migrant) {
	migrant <- best
}

// type OrganismAndBias struct {
// 	Bias float64
// 	Org  Organism
// }

type Migrant struct {
	Attractiveness float64
	Bias           float64
	Org            Organism
	PrevFitness    []float64
}

func MakeMigrant(bias float64, o Organism, pf []float64, nf []float64) *Migrant {
	// The migrant is initialized after the best organism from the generation
	// is picked. After,
	return &Migrant{
		Attractiveness: 0.0,
		Bias:           bias,
		Org:            o,
		PrevFitness:    pf,
	}
}

// Ai = Ai_prev + (n_pop + n_mig)
func (m *Migrant) CalculateAttractiveness(newFitness []float64) {
	popA := populationAttractiveness(m, newFitness)

	migA := migrantAttractiveness(m)

	m.Attractiveness += (popA + migA)
}

func populationAttractiveness(m *Migrant, newFitness []float64) float64 {
	residuals := 0.0
	for i := range newFitness {
		residuals += (m.PrevFitness[i] - newFitness[i])
	}

	npop := residuals / float64(len(newFitness))

	return math.Abs(npop)
}

func migrantAttractiveness(m *Migrant) float64 {
	// Since we only send one migrant, the nMig is 1.
	// The first migrant in the prev fitness is the most fit.
	nmig := m.PrevFitness[0] - m.Org.Fitness

	return math.Abs(nmig)
}

func MakeInitialPrevFitness(n int) []float64 {
	pf := make([]float64, n)

	for i := range pf {
		pf[i] = 99999.99
	}

	return pf
}

func GatherFitness(pop []Organism) []float64 {
	nf := make([]float64, len(pop))

	for i := range pop {
		nf[i] = pop[i].Fitness
	}

	return nf
}
