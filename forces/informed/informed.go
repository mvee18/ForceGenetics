package informed

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/guess"
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
var mu sync.Mutex

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

	mu.Lock()
	qt := r1.Intn(4)
	mu.Unlock()

	// fmt.Println(qt)

	if *flags.InitialGuess == "" {
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
					mu.Lock()
					chromosome[j] = utils.RandValueDomain(i + 2)
					mu.Unlock()
					if utils.RandBool() {
						chromosome[j] = -chromosome[j]
					}
				}
			}

			organism.DNA = append(organism.DNA, chromosome)
		}
	} else {
		organism.DNA = guess.MockB3LYP(*flags.InitialGuess)
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

func DirectedMutation(i InformedOrganism, g func(inf InformedOrganism) float64) InformedOrganism {
	previousFitness := i.Fitness
	// fmt.Println("previous fitness: ", previousFitness)
	iFinal := i

	// We must reevaluate the cost function at EVERY mutation.
	// This will take much longer, but it should give us better convergance over time...
	var wg sync.WaitGroup
	sema := make(chan struct{}, 4)

	for ind, v := range i.DNA {
		for j, gene := range v {
			iCopy := i
			sema <- struct{}{}
			wg.Add(1)
			go func(ind, j int, gene float64) {
				defer func() {
					<-sema
					wg.Done()
				}()

				mu.Lock()
				mutationChance := r1.Float64()
				mu.Unlock()

				mu.Lock()
				deltaNorm := r1.Float64()
				mu.Unlock()

				// If the corresponding direction is true (up), then add
				// the delta.
				// fmt.Printf("mutation chance is %v\n", *flags.MutationRateInformed)
				if mutationChance < *flags.MutationRateInformed {
					if i.Direction[ind][j] {
						// fmt.Println(deltaNorm, iCopy.DNA[ind][j])
						// The DNA at that index should have the denormalized delta added.
						newVal := (gene + deltaNorm*utils.SelectDomain(ind+2))
						incrementAndCheck(&newVal, ind+2)

						iFinal.DNA[ind][j] = newVal
						// fmt.Printf("on indices %d %d\n", ind, j)

						fitness := g(iCopy)

						// fmt.Printf("up: the old fitness is %v, the new fitness is %v\n", previousFitness, i.Fitness)

						// fmt.Println(fitness)

						// If the new fitness is worse than the old one, swap the direction.

						if fitness > previousFitness {
							// fmt.Println("fitness: ", fitness)
							iFinal.Direction[ind][j] = false
						}

					} else {
						// fmt.Println(deltaNorm, iCopy.DNA[ind][j])
						// If the mutation is down, then subtract.
						newVal := (iCopy.DNA[ind][j] - deltaNorm*utils.SelectDomain(ind+2))
						incrementAndCheck(&newVal, ind+2)

						iFinal.DNA[ind][j] = newVal
						// fmt.Printf("on indices %d %d\n", ind, j)

						fitness := g(iCopy)

						// fmt.Printf("down: the old fitness is %v, the new fitness is %v\n", previousFitness, i.Fitness)

						// If the new fitness is worse than the old one, swap the direction.
						if fitness > previousFitness {
							iFinal.Direction[ind][j] = true
						}
					}
				}
			}(ind, j, gene)

		}

	}

	wg.Wait()

	return iFinal
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
			mu.Lock()
			*v -= r1.Float64() * dom
			mu.Unlock()

		} else {
			// Else, if negative, add until within the bounds.
			mu.Lock()
			*v += r1.Float64() * dom
			mu.Unlock()
		}
	}

}

func GetBest(population InformedPopulation) InformedOrganism {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	return population[0]
}

func RunInformedGA(migrant chan models.Migrant) {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := CreateInformedPopulation()

	found := false
	generation := 0
	prevFitness := models.MakeInitialPrevFitness(len(population))

	for !found {
		generation++
		// fmt.Printf("Generation: %v\n", generation)
		bestOrganism := GetBest(population)

		nf := population.GatherFitness()

		mig := models.MakeMigrant(
			population.CalculateBias(),
			bestOrganism.Organism,
			prevFitness,
			nf,
		)

		models.AddMigrant(migrant, *mig)

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

			prevFitness = nf

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

func (p *InformedPopulation) AddImmigrant(migrant <-chan models.Migrant) {
	// Take the last organism (least fit) off.
	*p = (*p)[0 : len(*p)-1]

	// Need to convert migrant to informed type.
	org := <-migrant

	*p = append(*p, makeAndSetOrganism(&org.Org))
}

func (p *InformedPopulation) GatherFitness() []float64 {
	nf := make([]float64, len(*p))

	for i := range *p {
		nf[i] = (*p)[i].Fitness
	}

	return nf
}
