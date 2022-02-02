package islands

import (
	"ga/forces/informed"
	"ga/forces/models"
	"ga/forces/pseudo"
	trad "ga/forces/traditional"
	"math"
)

// Not exactly sure about the methods they share in common.
// Might want to add a TPopulation for the TGA.
type Island interface {
}

// We will gather than on the mig channel, then disperse them across the imm channels for each one.
func RunIslands(immPGA, immTGA, immIGA chan models.OrganismAndBias) {
	go pseudo.RunPGA(immPGA)

	go trad.RunTGA(immTGA)

	go informed.RunInformedGA(immIGA)

	migrantPool := make([]models.Organism, 0)

	for {
		select {
		case p := <-immPGA:
			// fmt.Printf("bias of pga: %v\n", p.Bias)
			migrantPool = append(migrantPool, p.Org)
			SendBestMigrant(p, immPGA, migrantPool)

		case t := <-immTGA:
			// fmt.Printf("bias of tga: %v\n", t.Bias)
			migrantPool = append(migrantPool, t.Org)
			SendBestMigrant(t, immTGA, migrantPool)

		case i := <-immIGA:
			// fmt.Printf("bias of iga: %v\n", i.Bias)
			migrantPool = append(migrantPool, i.Org)
			SendBestMigrant(i, immIGA, migrantPool)
		}
	}
}

func SendBestMigrant(o models.OrganismAndBias, mig chan<- models.OrganismAndBias, pool []models.Organism) {
	// First, we need to check if the pool has more than one member.
	if o.Bias >= 0.50 {
		if len(pool) > 1 {
			bestIndex, bestHD := 0, 0.0
			for i, v := range pool {
				hd := models.CalculateHD(o.Org, v)
				if hd > bestHD {
					bestIndex = i
					bestHD = hd
				}
			}

			mig <- models.OrganismAndBias{
				Org:  pool[bestIndex],
				Bias: 0.0,
			}

			RemoveIndex(pool, bestIndex)
		}

	} else {
		mig <- models.OrganismAndBias{
			Org:  o.Org,
			Bias: 0.0,
		}

	}
}

func RemoveIndex(s []models.Organism, index int) []models.Organism {
	return append(s[:index], s[index+1:]...)
}

// Ai = Ai_prev + (n_pop + n_mig)
func CalculateAttractiveness(m *models.Migrant, newFitness []float64) {
	popA := populationAttractiveness(m, newFitness)

	migA := migrantAttractiveness(m)

	m.Attractiveness += (popA + migA)
}

func populationAttractiveness(m *models.Migrant, newFitness []float64) float64 {
	residuals := 0.0
	for i := range newFitness {
		residuals += (m.PrevFitness[i] - newFitness[i])
	}

	npop := residuals / float64(len(newFitness))

	return math.Abs(npop)
}

func migrantAttractiveness(m *models.Migrant) float64 {
	// Since we only send one migrant, the nMig is 1.
	// The first migrant in the prev fitness is the most fit.
	nmig := m.PrevFitness[0] - m.Mig.Fitness

	return math.Abs(nmig)
}
