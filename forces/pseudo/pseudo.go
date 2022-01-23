package pseudo

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/selection"
	"ga/forces/utils"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

var r1 *rand.Rand

var (
	pseudoOutFile string = utils.NewOutputFile("pseudo/pseudo.out")
	bestPathFinal string = utils.NewOutputFile("pseudo/best/final")
)

func init() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(s1)
}

// If we declare new types here, we can keep the psuedo algorithm's methods
// separate from the TGA's methods. We can probably leverage this for a nice
// interface later. I also believe it will be cleaner to call methods on the
// organisms rather than functions?
type PPopulation []models.Organism

func CreatePseudoPopulation() (population PPopulation) {
	var wg sync.WaitGroup
	population = make(PPopulation, *flags.PopSize)

	sema := make(chan struct{}, 4)

	for i := 0; i < *flags.PopSize; i += 2 {
		sema <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-sema
				wg.Done()
			}()
			org := selection.CreateOrganism(*flags.NumAtoms)

			orgComp := complimentaryOrganism(&org)

			orgComp.SaveToFile(*flags.NumAtoms)
			orgComp.CalcFitness()

			population[i] = org
			population[i+1] = orgComp
		}(i)
	}
	wg.Wait()

	return
}

func complimentaryOrganism(o *models.Organism) models.Organism {
	comp := models.Organism{}

	newDna := o.DNA

	for i, v := range newDna {
		for j := range v {
			newDna[i][j] = -newDna[i][j]
		}
	}

	comp.DNA = newDna

	comp.SaveToFile(*flags.NumAtoms)

	return models.Organism(comp)
}

// We will now use the method prescribed in Chen, Zhong, Zhang (2010.) We cross
// over the adjacent (complimentary) parents and both offspring are added into
// the new population, replacing the old parents. Note that no new genes are
// added as there will be no mutation or natural selection.
func psuedoCrossover(pop PPopulation) PPopulation {
	newPop := make(PPopulation, len(pop))

	for i := 0; i < len(pop); i += 2 {

		p1 := (pop)[i]
		p2 := (pop)[i+1]

		c1, c2 := performPseudoCrossover(p1, p2)

		c1.SaveToFile(*flags.NumAtoms)
		c2.SaveToFile(*flags.NumAtoms)

		c1.CalcFitness()
		c2.CalcFitness()

		newPop[i] = c1
		newPop[i+1] = c2
	}

	return newPop
}

// The parents are the same exact size always.
func performPseudoCrossover(p1, p2 models.Organism) (models.Organism, models.Organism) {
	c1 := models.Organism{
		DNA:     duplicateDNA(p1),
		Path:    "",
		Fitness: 0.0,
	}

	c2 := models.Organism{
		DNA:     duplicateDNA(p2),
		Path:    "",
		Fitness: 0.0,
	}

	for i, v := range p1.DNA {
		for j := range v {
			if crossoverChance() {

				// fmt.Printf("p1 gene: %v\n, p2 gene: %v\n", p1.DNA[i][j], p2.DNA[i][j])

				original := p1.DNA[i][j]
				// gene is the original value of p1[i][j].
				c1.DNA[i][j] = p2.DNA[i][j]
				c2.DNA[i][j] = original

				// fmt.Printf("after::: p1 gene: %v\n, p2 gene: %v\n", p1.DNA[i][j], p2.DNA[i][j])

			} else {
				continue
			}
		}
	}

	return c1, c2
}

func crossoverChance() bool {
	chance := r1.Float64()
	if chance <= *flags.PseudoCrossOverRate {
		return true
	}

	return false
}

func duplicateDNA(p1 models.Organism) models.DNA {
	duplicate := make(models.DNA, len(p1.DNA))
	for i := range p1.DNA {
		duplicate[i] = make(models.Chromosome, len(p1.DNA[i]))
		copy(duplicate[i], p1.DNA[i])
	}

	return duplicate
}

func delFolders(o []models.Organism, topOrganism models.Organism) {
	for _, v := range o {
		if v.Path == topOrganism.Path {
			continue
		} else {
			os.RemoveAll(path.Dir(v.Path))
		}
	}
}

func (p *PPopulation) AddImmigrant(migrant <-chan models.Organism) {
	// Take the last organism (least fit) off.
	*p = (*p)[0 : len(*p)-1]

	*p = append(*p, <-migrant)
}

func RunPGA(migrant chan models.Organism) {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := CreatePseudoPopulation()

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := selection.GetBest(population)

		// fmt.Println("added migrant from pga.")
		models.AddMigrant(migrant, bestOrganism)

		if bestOrganism.Fitness < *flags.FitnessLimit {
			found = true

			err := bestOrganism.LogFinalOrganism(start, pseudoOutFile, bestPathFinal)
			if err != nil {
				log.Fatalln(err)
			}

			return

		} else {
			population = psuedoCrossover(population)

			if generation != 0 {
				population.AddImmigrant(migrant)
			}

			if generation%10 == 0 {
				bestPath := utils.NewOutputFile(fmt.Sprintf("pseudo/best/%d", generation))
				err := bestOrganism.LogIntermediateOrganism(generation, start, pseudoOutFile, bestPath)
				if err != nil {
					log.Fatalln(err)
				}

				if generation >= *flags.GenLimit {
					models.LogTerminated(pseudoOutFile)
				}
			}

			delFolders(population, bestOrganism)
		}

	}

}
