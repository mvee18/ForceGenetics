package trad

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
	"time"
)

var (
	tradOutFile   string = utils.NewOutputFile("traditional/trad.out")
	bestPathFinal string = utils.NewOutputFile("traditional/best/final")
)

func RunTGA(migrant chan models.OrganismAndBias) {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := selection.CreatePopulation()

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := selection.GetBest(population)

		// Add migrant to pool.
		// fmt.Println("migrant added from traditional ga")
		best := models.OrganismAndBias{
			Org:  bestOrganism,
			Bias: models.CalculateBias(population),
		}

		models.AddMigrant(migrant, best)

		if bestOrganism.Fitness < *flags.FitnessLimit {
			found = true

			bestOrganism.LogFinalOrganism(start, tradOutFile, bestPathFinal)

			return

		} else {
			pool := selection.CreatePool(population, models.TargetFrequencies)

			// If the generation is 0, there won't be any immigrants
			// to put in the pool from the channel.
			if generation != 0 {
				models.AddImmigrant(&pool, migrant)
			}

			population = selection.NaturalSelection(pool, population, models.TargetFrequencies)

			if generation%10 == 0 {
				bestPath := utils.NewOutputFile(fmt.Sprintf("traditional/best/%d", generation))
				err := bestOrganism.LogIntermediateOrganism(generation, start, tradOutFile, bestPath)
				if err != nil {
					log.Fatalln(err)
				}

				if generation >= *flags.GenLimit {
					models.LogTerminated(tradOutFile)
				}
			}

			delFolders(pool, bestOrganism)
			delFolders(population, bestOrganism)
		}

	}

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
