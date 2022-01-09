package informed

import (
	"ga/forces/flags"
	"ga/forces/quadratic"
	"math/rand"
	"sort"
)

func CreatePool(population InformedPopulation, target []float64) (pool InformedPopulation) {
	pool = make(InformedPopulation, 0)

	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	fraction := *flags.PoolSize * float64(*flags.PopSize)

	top := population[0:int(fraction)]

	// We take the top 50% of organisms, as prescribed in Haupt and Haupt p. 54.
	pool = append(pool, top...)

	//	fmt.Println("The length of the pool is", len(pool))

	return
}

func NaturalSelection(pool InformedPopulation, population InformedPopulation, target []float64) InformedPopulation {
	// Children = [pool + empty slice]; len = population.
	next := make(InformedPopulation, len(population)-len(pool))

	children := append(pool, next...)

	//fmt.Println("The original length of next is: ", len(children))

	//fmt.Printf("Length of pop minus pool is %v\n", len(population)-len(pool))

	// Remember the principle of Independent Assortment. Each chromosome should go through crossover individually.
	for i := len(pool); i < len(population); i++ {
		/*
			r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
			a := pool[r1]
			b := pool[r2]
		*/

		// We could do some boolean mixing...?
		a, b, c := tournamentRound(pool)

		child := quadratic.QuadraticTerms(&a.Organism, &b.Organism, &c.Organism)

		err := child.SaveToFile(*flags.NumAtoms)
		if err != nil {
			panic(err)
		}

		// Make into informed organism.
		informed := InformedOrganism{
			Organism:  child,
			Direction: [3][]bool{},
		}

		// The organism does not have a new fitness. This is necessary
		// for the directed mutation.
		informed.CalcFitness()

		// We combined the parents to yield a new velocity.
		informed.CombinedVelocity(a.Direction, b.Direction, c.Direction)

		DirectedMutation(&informed, (*InformedOrganism).CalcFitness)

		children[i] = informed
	}

	//fmt.Println("The length of next is: ", len(children))
	return children
}

func tournamentRound(pool InformedPopulation) (InformedOrganism, InformedOrganism, InformedOrganism) {
	// This grabs three random organisms from the pool and sorts them.
	makeBracket := func() InformedPopulation {
		round := make(InformedPopulation, 0)

		for i := 0; i < *flags.TournamentPool; i++ {
			index := rand.Intn(len(pool))
			round = append(round, pool[index])
		}
		// We only need to sort the round once.
		sort.SliceStable(round, func(i, j int) bool {
			return round[i].Fitness < round[j].Fitness
		})

		return round
	}

	// This populates the slices a, b, which will be used in the next step.
	a, b, c := makeBracket(), makeBracket(), makeBracket()

	// Since these slices are sorted, the first element is the most fit organism.
	return a[0], b[0], c[0]
}
