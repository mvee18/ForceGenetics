package selection

import (
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"ga/forces/quadratic"
	"ga/forces/utils"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
)

// MutationRate is the rate of mutation
// var MutationRate = 0.0004

// PopSize is the size of the population (was 500)
// var PopSize = 600

// The number of atoms is critical to the size of the organism.
// var NumAtoms = 6

// Size of breeding pool (was 30)
// var PoolSize = 35

// var FitnessLimit = 1.0

/*
var TargetFrequencies = []float64{
	//820.24, 804.08, 737.75,
	//580.87, 573.06, 525.60,
	//363.23, 274.64, 194.26,
	//170.33, 170.26, 0.34,
	//0.00, 0.00, 0.00,
	//0.00, 0.00, 0.00,
	3943.98, 3833.99, 1651.33,
	0.02, 0.00, 0.00,
	0.00, 0.00, 0.00,
}
*/

func CreateOrganism(numAtoms int) (organism models.Organism) {
	// This iterates over the derivative levels to fill in the DNA for each
	// organisms on the 3 chromosomes.
	// Ex. if d = 3 ==> 2, then we get Ch 0, 1 filled.
	organism = models.Organism{
		DNA:     []models.Chromosome{},
		Path:    "",
		Fitness: 0,
	}

	for i := 0; i < *flags.DerivativeLevel-1; i++ {
		organismSize := utils.GetNumForceConstants(numAtoms, i+2)
		chromosome := make([]float64, organismSize)
		for j := 0; j < organismSize; j++ {
			chromosome[j] = (0.0 + rand.Float64()*(*flags.Domain-0.0))
			if utils.RandBool() {
				chromosome[j] = -chromosome[j]
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

// Dir is the temp directory where the file will be stored.
func copyFile(dir *string, fortFile *string, d *int) error {
	// New variable. We don't want to accidentally edit the derivative level.
	derivativeIndex := *d - 2

	var fortFiles = []string{"fort.15", "fort.30", "fort.40"}
	for i := 0; i <= derivativeIndex; i++ {
		*fortFile = fortFiles[i]
		// That is, if this is the last loop.
		if i == derivativeIndex {
			break
		} else {
			// Otherwise, we need to copy the files below that derivative
			// level.
			outputFile := path.Join(*dir, fortFiles[i])
			_, err := copy(*fortFile, outputFile)
			if err != nil {
				fmt.Printf("error copying file, %v\n", err)
				return err
			}

		}

	}

	return nil
}

// Residual Sum of Squares

func CreatePopulation() (population []models.Organism) {
	var wg sync.WaitGroup
	population = make([]models.Organism, *flags.PopSize)

	sema := make(chan struct{}, 4)

	for i := 0; i < *flags.PopSize; i++ {
		sema <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-sema
				wg.Done()
			}()
			org := CreateOrganism(*flags.NumAtoms)
			population[i] = org
		}(i)
	}
	wg.Wait()

	return
}

func CreatePool(population []models.Organism, target []float64) (pool []models.Organism) {
	pool = make([]models.Organism, 0)
	// get top best fitting organisms
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	// fmt.Printf("The top fitness of the pool is: %v %v %v", population[0].Fitness, population[1].Fitness, population[2].Fitness)

	// This is what fraction survives to the next generation.
	//	fmt.Println("length population, ", len(population))
	fraction := *flags.PoolSize * float64(*flags.PopSize)

	top := population[0:int(fraction)]

	//fmt.Println("The top fitness is, and the least is:", top[0].Fitness, top[len(top)-1].Fitness)
	//	fmt.Println("Path to top is: ", top[0].Path)

	// bottom := population[*PoolSize+2:]

	// if there is no difference between the top  organisms, the population is stable
	// and we can't get generate a proper breeding pool so we make the pool equal to the
	// population and reproduce the next generation

	// This might be necessary?
	/*
		if top[len(top)-1].Fitness-top[0].Fitness == 0 {
			pool = population
			return
		}
	*/

	// create a pool for next generation
	/*
		for i := 0; i < len(top)-1; i++ {
			num := (top[*PoolSize].Fitness - top[i].Fitness)
			fmt.Println("num: ", num)
			for n := int64(0); n < int64(num); n++ {
				pool = append(pool, top[i])
			}
		}
	*/

	// We take the top 50% of organisms, as prescribed in Haupt and Haupt p. 54.
	pool = append(pool, top...)

	//	fmt.Println("The length of the pool is", len(pool))

	return
}

// perform natural selection to create the next generation
// We should use a weighted method.
func NaturalSelection(pool []models.Organism, population []models.Organism, target []float64) []models.Organism {
	// Children = [pool + empty slice]; len = population.
	next := make([]models.Organism, len(population)-len(pool))

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

		a, b, c := tournamentRound(pool)

		child := quadratic.QuadraticTerms(&a, &b, &c)
		child.Mutate()

		err := child.SaveToFile(*flags.NumAtoms)
		if err != nil {
			panic(err)
		}

		child.CalcFitness()

		children[i] = child
	}

	//fmt.Println("The length of next is: ", len(children))
	return children
}

func tournamentRound(pool []models.Organism) (models.Organism, models.Organism, models.Organism) {
	// This grabs three random organisms from the pool and sorts them.
	makeBracket := func() []models.Organism {
		round := make([]models.Organism, 0)

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

// crosses over 2 Organism strings
// We should use the blending method prescribed in Haupt and Haupt.
// Furthermore, each cross should produce two offspring.
/*
func crossover(d1 models.Organism, d2 models.Organism) models.Organism {
	childDNA := make(DNA, len(d1.DNA))

	for i, pChr := range d1.DNA {
		ch := make(models.Chromosome, len(pChr))
		childDNA[i] = ch
	}

	child := Organism{
		DNA:     childDNA,
		Path:    "",
		Fitness: 0,
	}

	// if rand.Float64() <= *CrossOverRate {
	// This is the simple method of

	for i, chr := range d1.DNA {
		mid := rand.Intn(len(chr))
		for j := 0; j < len(chr); j++ {
			if j < mid {
				child.DNA[i][j] = d1.DNA[i][j]
				//			} else if i > mid {
				//		child.DNA[i][j] = d2.DNA[i][j]
			} else if i >= mid {
				if RandBool() {
					child.DNA[i][j] = crossOverA(d1.DNA[i][j], d2.DNA[i][j])
				} else {
					child.DNA[i][j] = crossOverB(d1.DNA[i][j], d2.DNA[i][j])
				}
			}
		}

	}

	return child
}

// Where m is the mother chromosome and d is the father chromosome.
func crossOverA(m float64, d float64) float64 {
	// CrossoverPoint
	beta := rand.Float64()
	pNew := m - beta*m + beta*d

	return pNew
}

func crossOverB(m float64, d float64) float64 {
	// CrossoverPoint
	beta := rand.Float64()
	pNew := m + beta*m - beta*d

	return pNew
}
*/

func GetBest(population []models.Organism) models.Organism {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	return population[0]
}

var OutputPath string

func init() {
	outputPath, err := filepath.Abs(*flags.OutFile)
	if err != nil {
		log.Fatalf("Error in getting output path, %v\n", err)
	}

	OutputPath = outputPath

	f, err := os.Create(OutputPath)
	if err != nil {
		log.Fatalf("Error generating output file, %v\n", err)
	}

	f.Close()

	if os.Stat("best"); !os.IsNotExist(err) {
		os.RemoveAll("best")
	}

	err = os.Mkdir("best", 0700)
	if err != nil {
		panic(err)
	}

	models.TargetFrequencies, models.TargetRotational, models.TargetFund = utils.ReadInput(*flags.FreqInputFile)
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

/*
func countOpenFiles() int64 {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	lines := strings.Split(string(out), "\n")
	return int64(len(lines) - 1)
}
*/
