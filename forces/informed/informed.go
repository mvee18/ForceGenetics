package informed

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/utils"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"sort"
	"sync"
	"time"
)

var r1 *rand.Rand

var (
	informedOutFile string = utils.NewOutputFile("informed/informed.out")
	bestPathFinal   string = utils.NewOutputFile("informed/best/final")
)

func init() {
	s1 := rand.NewSource(time.Now().UnixNano() + 2561)
	r1 = rand.New(s1)
}

// Gozali et. al really need to publish their code because their description of
// the implementations make no sense.

type InformedPopulation []InformedOrganism

// We need to add the direction to comply with the simplified swarm particle mode
// that Gozali suggests. As a result, we need to redefine the methods.
type InformedOrganism struct {
	models.Organism
	Direction [3][]bool
}

func makeAndSetOrganism(org *models.Organism) InformedOrganism {
	iOrg := InformedOrganism{
		Organism:  *org,
		Direction: [3][]bool{},
	}

	iOrg.CreateVelocity()

	return iOrg
}

// TODO: Split the population in half with quartile/random. IDK make it look
// clean somehow without tons of paramter passing.
func CreateInformedPopulation() (population InformedPopulation) {
	var wg sync.WaitGroup
	population = make([]InformedOrganism, *flags.PopSize)

	sema := make(chan struct{}, 4)

	for i := 0; i < *flags.PopSize; i++ {
		sema <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-sema
				wg.Done()
			}()

			if i <= *flags.PopSize {
				org := CreateInformedOrganism(*flags.NumAtoms, true)
				population[i] = makeAndSetOrganism(&org)

			} else {
				org := CreateInformedOrganism(*flags.NumAtoms, false)
				population[i] = makeAndSetOrganism(&org)
			}
		}(i)
	}
	wg.Wait()

	return
}

func (p *InformedOrganism) CreateVelocity() {
	var v [3][]bool

	// This should allocate the second dimension to be the same length as
	// the DNA... a 1:1 mapping for velocity on the same index.
	// Thus, the corresponding velocity at DNA[i][j] is V[i][j].
	for i, val := range p.DNA {
		v[i] = make([]bool, len(val))
	}

	p.Direction = v

	for j, k := range p.Direction {
		for l := range k {
			p.Direction[j][l] = (utils.RandBool())
		}
	}
}

func (p *InformedOrganism) CombinedVelocity(a, b, c [3][]bool) {
	var v [3][]bool

	for i, val := range a {
		v[i] = make([]bool, len(val))
	}

	for i, val := range a {
		for j := range val {
			v[i][j] = atLeastTwo(a[i][j], b[i][j], c[i][j])
		}
	}

	p.Direction = v
}

// Returns whether true or false is more common.
func atLeastTwo(a, b, c bool) bool {
	return a && (b || c) || (b && c)
}

func CreateInformedOrganism(numAtoms int, quartile bool) (organism models.Organism) {
	// This iterates over the derivative levels to fill in the DNA for each
	// organisms on the 3 chromosomes.
	// Ex. if d = 3 ==> 2, then we get Ch 0, 1 filled.
	organism = models.Organism{
		DNA:     []models.Chromosome{},
		Path:    "",
		Fitness: 0,
	}

	// The organisms needs its quartile determined ahead of time so all of
	// its DNA in the same quartile frame.

	qt := r1.Intn(4)

	// fmt.Println(qt)

	for i := 0; i < *flags.DerivativeLevel-1; i++ {
		organismSize := utils.GetNumForceConstants(numAtoms, i+2)
		chromosome := make([]float64, organismSize)

		qv := QuartileValues(i + 2)

		for j := 0; j < organismSize; j++ {
			if quartile {
				chromosome[j] = QuartileValueDomain(qv, qt)
				if utils.RandBool() {
					chromosome[j] = -chromosome[j]
				}

			} else {
				chromosome[j] = utils.RandValueDomain(i + 2)
				if utils.RandBool() {
					chromosome[j] = -chromosome[j]
				}
			}
		}

		organism.DNA = append(organism.DNA, chromosome)
	}

	err := organism.SaveToFile(numAtoms)
	if err != nil {
		fmt.Printf("Error in saving organism to file, %v\n", err)
	}

	organism.CalcFitness()

	return organism
}

func QuartileValueDomain(v [4]float64, qt int) float64 {
	// v := [4]float64{q1, q2, q3, ub}

	// Q1 is within first 25%, Q2 within first half, Q2 first 75%.
	quartile := v[qt]

	return (0.0 + rand.Float64()*(quartile-0.0))
}

// As prescribed in the Gozali paper, we need to follow the IGA steps.
// We first initialize a quartile system for the lower/upper bound.
func QuartileValues(dn int) [4]float64 {
	ub := utils.SelectDomain(dn)
	lb := 0.0

	q2 := (lb + ub) / 2
	q1 := q2 * 0.5
	q3 := q2 * 1.5

	return [4]float64{q1, q2, q3, ub}
}

func DirectedMutation(i *InformedOrganism, g func(inf *InformedOrganism)) {
	previousFitness := i.Fitness

	// We must reevaluate the cost function at EVERY mutation.
	// This will take much longer, but it should give us better convergance over time...
	for ind, v := range i.DNA {
		for j := range v {
			mutationChance := r1.Float64()

			deltaNorm := r1.Float64()

			// If the corresponding direction is true (up), then add
			// the delta.
			if mutationChance < *flags.MutationRateInformed {
				if i.Direction[ind][j] {
					// The DNA at that index should have the denormalized delta added.
					i.DNA[ind][j] += deltaNorm * utils.SelectDomain(ind+2)
					incrementAndCheck(&i.DNA[ind][j], ind+2)

					g(i)

					// fmt.Printf("up: the old fitness is %v, the new fitness is %v\n", previousFitness, i.Fitness)

					// If the new fitness is worse than the old one, swap the direction.
					if i.Fitness > previousFitness {
						i.Direction[ind][j] = false
					}

				} else {
					// If the mutation is down, then subtract.
					i.DNA[ind][j] -= deltaNorm * utils.SelectDomain(ind+2)
					incrementAndCheck(&i.DNA[ind][j], ind+2)

					g(i)

					// fmt.Printf("down: the old fitness is %v, the new fitness is %v\n", previousFitness, i.Fitness)

					// If the new fitness is worse than the old one, swap the direction.
					if i.Fitness > previousFitness {
						i.Direction[ind][j] = true
					}
				}
			}
		}
	}

}

// This function ensures that the bounds aren't exceeded by the directed
// mutation. It randomly adds/subtracts values from the gene until it is within
// the bounds.
func incrementAndCheck(v *float64, dn int) {
	dom := utils.SelectDomain(dn)

	// If the value is out of bounds...
	for math.Abs(*v) > dom {
		if *v > 0 {
			// If positive, subtract until it is below...
			*v -= r1.Float64() * dom

		} else {
			// Else, if negative, add until within the bounds.
			*v += r1.Float64() * dom
		}
	}

}

func GetBest(population InformedPopulation) InformedOrganism {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	return population[0]
}

func RunInformedGA(migrant chan models.OrganismAndBias) {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := CreateInformedPopulation()

	found := false
	generation := 0
	for !found {
		generation++
		// fmt.Printf("Generation: %v\n", generation)
		bestOrganism := GetBest(population)

		best := models.OrganismAndBias{
			Org:  bestOrganism.Organism,
			Bias: population.CalculateBias(),
		}

		models.AddMigrant(migrant, best)

		if bestOrganism.Fitness < *flags.FitnessLimit {
			found = true

			bestOrganism.LogFinalOrganism(start, informedOutFile, bestPathFinal)

			return

		} else {
			pool := CreatePool(population, models.TargetFrequencies)

			if generation != 0 {
				population.AddImmigrant(migrant)
				// fmt.Println("iga received migrant and is continuing")
			}

			population = NaturalSelection(pool, population, models.TargetFrequencies)

			if generation%10 == 0 {
				bestPath := utils.NewOutputFile(fmt.Sprintf("informed/best/%d", generation))
				err := bestOrganism.LogIntermediateOrganism(generation, start, informedOutFile, bestPath)
				if err != nil {
					log.Fatalln(err)
				}

				if generation >= *flags.GenLimit {
					models.LogTerminated(informedOutFile)
				}
			}

			delFolders(pool, bestOrganism)
			delFolders(population, bestOrganism)

		}

	}

}

func (i InformedPopulation) CalculateBias() float64 {
	pop := make([]models.Organism, len(i))
	for ind, v := range i {
		pop[ind] = v.Organism
	}

	return models.CalculateBias(pop)
}

func delFolders(o InformedPopulation, topOrganism InformedOrganism) {
	for _, v := range o {
		if v.Path == topOrganism.Path {
			continue
		} else {
			os.RemoveAll(path.Dir(v.Path))
		}
	}
}

func (p *InformedPopulation) AddImmigrant(migrant <-chan models.OrganismAndBias) {
	// Take the last organism (least fit) off.
	*p = (*p)[0 : len(*p)-1]

	// Need to convert migrant to informed type.
	org := <-migrant

	*p = append(*p, makeAndSetOrganism(&org.Org))
}
