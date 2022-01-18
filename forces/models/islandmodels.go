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

func CalculateHD(ub float64, a, b Organism) {
	totalHD := 0.0

	for i, chr := range a.DNA {
		chrHD := 0.0

		for j := range chr {
			chrHD += ((a.DNA[i][j] + ub) / (2 * ub)) - ((b.DNA[i][j] + ub) / (2 * ub))
		}

		totalHD += math.Abs(chrHD)
	}
}
