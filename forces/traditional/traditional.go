package trad

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/selection"
	"math/rand"
	"os"
	"path"
	"time"
)

func RunTGA(immigrant <-chan models.Organism, migrant chan<- models.Organism) {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := selection.CreatePopulation()

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := selection.GetBest(population)

		// Add migrant to pool.
		go AddMigrant(migrant, bestOrganism)

		if bestOrganism.Fitness < *flags.FitnessLimit {
			found = true

			f, err := os.OpenFile(selection.OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
			pool := selection.CreatePool(population, models.TargetFrequencies)

			// If the generation is 0, there won't be any immigrants
			// to put in the pool from the channel.
			if generation != 0 {
				AddImmigrant(&pool, immigrant)
			}

			population = selection.NaturalSelection(pool, population, models.TargetFrequencies)

			if generation%10 == 0 {
				sofar := time.Since(start)

				summaryStep := fmt.Sprintf("The path to the best organism is %v.\n \nTime taken so far: %s | generation: %d | fitness: %f | pool size: %d\n", bestOrganism.Path, sofar, generation, bestOrganism.Fitness, len(pool))

				f, err := os.OpenFile(selection.OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

			delFolders(pool, bestOrganism)
			delFolders(population, bestOrganism)
		}

	}

}

func AddImmigrant(p *[]models.Organism, migrant <-chan models.Organism) {
	// Take the last organism (least fit) off.
	*p = (*p)[0 : len(*p)-1]

	*p = append(*p, <-migrant)
}

func AddMigrant(migrant chan<- models.Organism, best models.Organism) {
	migrant <- best
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
