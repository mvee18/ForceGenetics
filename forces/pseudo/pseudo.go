package pseudo

import (
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/selection"
	"sync"
)

func CreatePseudoPopulation() (population []models.Organism) {
	var wg sync.WaitGroup
	population = make([]models.Organism, *flags.PopSize)

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

	return comp
}
