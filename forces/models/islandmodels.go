package models

// Unfortunately, we have multiple chromosomes.
// As a result, we have to modify the version presented by Tsutshi, 1995.
func CalculateBias(p []Organism) {

	// The formula specifies we compare the same gene across all the organisms
	// in the population first, not each organism.
	n := len(p)

	for _, v := range p {
		for j, chromosome := range v.DNA {
			// We need all of the chromosomes in this column.
			lChr := len(chromosome)
			for k := range chromosome {
				// Iterate again across all of the organisms (rows).
				for m := range p {
					c := 0.0
					biasValue(&c, p[m].DNA[j][k], n, lChr)
				}

			}
		}
	}

}

func biasValue(count *float64, gene float64, n, l int) {

}
