package pseudo

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/selection"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"
)

var r1 *rand.Rand

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

func runPGA() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := CreatePseudoPopulation()

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := selection.GetBest(population)

		if bestOrganism.Fitness < *flags.FitnessLimit {
			found = true

			f, err := os.OpenFile(*flags.OutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			foundString := fmt.Sprintf("The path to the best organism is %v\n", bestOrganism.Path)

			if _, err = f.WriteString(foundString); err != nil {
				panic(err)
			}

			if _, err = f.WriteString("Yes, the superior fighter is clear. Succcessful termination.\n"); err != nil {
				panic(err)
			}

			f.Close()

			bestPath := "best/final"
			bestErr := bestOrganism.SaveBestOrganism(*flags.NumAtoms, bestPath)
			if bestErr != nil {
				fmt.Printf("Error saving best organism, %v\n", err)
			}

			elapsed := time.Since(start)
			fmt.Printf("\nTotal time taken: %s\n", elapsed)

			return

		} else {
			population = psuedoCrossover(population)

			if generation%10 == 0 {
				sofar := time.Since(start)

				summaryStep := fmt.Sprintf("The path to the best organism is %v.\n \nTime taken so far: %s | generation: %d | fitness: %f", bestOrganism.Path, sofar, generation, bestOrganism.Fitness)

				f, err := os.OpenFile(*flags.OutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}

				if _, err = f.WriteString(summaryStep); err != nil {
					panic(err)
				}

				f.Close()

				bestPath := fmt.Sprintf("best/%d", generation)
				bestErr := bestOrganism.SaveBestOrganism(*flags.NumAtoms, bestPath)
				if bestErr != nil {
					fmt.Printf("Error saving best organism, %v\n", err)
				}

				if generation >= *flags.GenLimit {
					f, err := os.OpenFile(selection.OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					if err != nil {
						panic(err)
					}

					if _, err = f.WriteString("Terminated. Maximum number of generations reached."); err != nil {
						panic(err)
					}

					f.Close()

					fmt.Println("Maximum number of generations reached.")
					os.Exit(0)
				}
			}

			delFolders(population, bestOrganism)
		}

	}

}
