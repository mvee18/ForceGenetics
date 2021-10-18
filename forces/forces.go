package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/ntBre/chemutils/summarize"
)

var ErrNullSummarize = errors.New("null value of organism")
var ErrCalcFitness = errors.New("error calculating fitness of organism")
var ErrCalcHarmFitness = errors.New("error calculating harm fitness")
var ErrCalcRotFitness = errors.New("error calculating rot fit")
var ErrCalcFundFitness = errors.New("error calculating fund fit")
var ErrNaNFitness = errors.New("Not a number fitness")

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

var FortFiles = []string{"fort.15", "fort.30", "fort.40"}

var TargetFrequencies = []float64{
	3943.98, 3833.99, 1651.33,
	0, 0, 0,
	0, 0, 0,
}

var TargetRotational = []float64{
	14.5054957,
	9.2636424,
	27.6557350,
}

var TargetFund = []float64{
	3753.156,
	3656.489,
	1598.834,
}

type DNA []Chromosome
type Chromosome []float64

// The organism is going to be the array of force constants.
// We should be able to represent this as a one dimensions array,
// then reconstruct it to run it in spectro.
type Organism struct {
	DNA     DNA
	Path    string
	Fitness float64
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func GetNumForceConstants(n int, dn int) int {
	switch dn {
	case 2:
		return int(math.Pow(float64(n), 2)) * 3 * 3

	case 3:
		c := 0
		for i := 0; i <= n*3-1; i++ {
			for j := 0; j <= i; j++ {
				for k := 0; k <= j; k++ {
					c++
				}
			}
		}
		return c

	case 4:
		c := 0
		for i := 0; i <= n*3-1; i++ {
			for j := 0; j <= i; j++ {
				for k := 0; k <= j; k++ {
					for l := 0; l <= k; l++ {
						c++
					}
				}
			}
		}
		return c

	default:
		panic("derivative level invalid, must be 2,3,4.")
	}

}

func CreateOrganism(numAtoms int) (organism Organism) {
	// This iterates over the derivative levels to fill in the DNA for each
	// organisms on the 3 chromosomes.
	// Ex. if d = 3 ==> 2, then we get Ch 0, 1 filled.
	organism = Organism{
		DNA:     []Chromosome{},
		Path:    "",
		Fitness: 0,
	}

	for i := 0; i < *DerivativeLevel-1; i++ {
		organismSize := GetNumForceConstants(numAtoms, i+2)
		chromosome := make([]float64, organismSize)
		for j := 0; j < organismSize; j++ {
			chromosome[j] = (rand.Float64())
			if RandBool() {
				chromosome[j] = -chromosome[j]
			}
		}

		organism.DNA = append(organism.DNA, chromosome)
	}

	err := organism.saveToFile(numAtoms)
	if err != nil {
		fmt.Printf("Error in saving organism to file, %v\n", err)
	}

	organism.calcFitness()

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

// Before we can calc the fitness, we have to save the files so that spectro can use them.
// Let's use temp files.
func (d *Organism) saveToFile(natoms int) error {
	dir, err := ioutil.TempDir(".", "forceOrganisms")
	if err != nil {
		log.Fatalf("Could not open temp dir, %v\n", err)
		return err
	}

	for i, chr := range d.DNA {
		fortFile := FortFiles[i]
		tempfn := path.Join(dir, fortFile)
		organismFile, err := os.Create(tempfn)
		if err != nil {
			log.Fatalf("Could not open temp file, %v\n", err)
			return err
		}

		fmt.Fprintf(organismFile, "%5d%5d", natoms, GetNumForceConstants(natoms, i+2))
		for j := range chr {
			if j%3 == 0 {
				fmt.Fprintf(organismFile, "\n")
			}
			fmt.Fprintf(organismFile, "%20.10f", d.DNA[i][j])
		}
		organismFile.Write([]byte("\n"))

		d.Path = organismFile.Name()

		organismFile.Close()
	}

	// Now we need to format the file correctly.
	// Spectro is 20.12f
	return nil
}

func (d *Organism) SaveBestOrganism(natoms int, filePath string) error {

	fmt.Printf("\nThe filePath is %s\n", filePath)

	err := os.MkdirAll(filePath, 0700)
	if err != nil {
		log.Fatalf("Could not open temp dir, %v\n", err)
		return err
	}

	// fmt.Printf("The best organism is %#v\n", d)

	for i, chr := range d.DNA {
		fortFile := FortFiles[i]
		tempfn := path.Join(filePath, fortFile)
		organismFile, err := os.Create(tempfn)
		if err != nil {
			log.Fatalf("Could not open temp file, %v\n", err)
			return err
		}

		fmt.Fprintf(organismFile, "%5d%5d", natoms, GetNumForceConstants(natoms, i+2))
		for j := range chr {
			if j%3 == 0 {
				fmt.Fprintf(organismFile, "\n")
			}
			fmt.Fprintf(organismFile, "%20.10f", d.DNA[i][j])
		}
		organismFile.Write([]byte("\n"))

		d.Path = organismFile.Name()

		organismFile.Close()
	}

	return nil
}

// To calculate the fitness, we must run it through spectro.
// We can save the results to a temp file and get a difference.
// The smaller the differences, the greater the fitness (1/difference).
func (d *Organism) calcFitness() {
	// Let's begin with opening the fort file in each organism.
	f, err := os.Open(d.Path)
	if err != nil {
		log.Fatalf("Error opening file of path, %v\n", err)
	}

	defer f.Close()

	input, err := os.Open(*PathToSpectroIn)
	if err != nil {
		log.Fatalf("Error opening file of spectro in, %v\n", err)
	}

	defer input.Close()

	spectroAbs, err := filepath.Abs(*PathToSpectro)
	if err != nil {
		log.Fatalf("Error finding abs path, %v\n", err)
	}

	cmd := exec.Cmd{
		Path:  spectroAbs,
		Stdin: input,
		Dir:   path.Dir(f.Name()),
	}

	outBytes, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running cmd, %v\n", err)
	}

	//	ioutil.WriteFile("test.out", outBytes, 0777)

	parseOutput(d, outBytes, *DerivativeLevel)
}

// TODO: Add a constraint on negative frequencies. Should reduce the fitness of
// the organism, even if the freq get closer.
func parseOutput(d *Organism, by []byte, derivative int) {
	r := bytes.NewReader(by)
	result := summarize.Spectro(r)

	// fmt.Printf("%#v", result)
	// fmt.Println(d.Path)

	fitness := 9999.0

	var err error

	switch derivative {
	case 2:
		fitness, err = calcDifference(result.Harm, TargetFrequencies)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Fatalln(ErrCalcHarmFitness)
			}
		}
	case 3:
		harmFitness, err := calcDifference(result.Harm, TargetFrequencies)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Fatalln(ErrCalcHarmFitness)
			}
		}

		rotFitness, err := calcDifference(result.Rots[0], TargetFrequencies)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Fatalln(ErrCalcRotFitness)
			}
		}

		fitness = harmFitness + rotFitness

	case 4:
		if len(result.Fund) == 0 || len(result.Harm) == 0 || len(result.Rots[0]) == 0 {
			// fmt.Printf("Singular matrix organism, %v\n", d.Path)
			d.Fitness = 99999.99
			return
		}

		harmFitness, err := calcDifference(result.Harm, TargetFrequencies)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Printf("%v in organism %v\n", err, d.Path)
				d.Fitness = 99999.99
				return
			}
		}

		rotFitness, err := calcDifference(result.Rots[0], TargetFrequencies)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Printf("%v in organism %v\n", err, d.Path)
				d.Fitness = 99999.99
				return
			}
		}

		fundFitness, err := calcDifference(result.Fund, TargetFund)
		if err != nil {
			if err == ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Printf("%v in organism %v\n", err, d.Path)
				d.Fitness = 99999.99
				return
			}
		}

		fitness = harmFitness + rotFitness + fundFitness
	}

	d.Fitness = fitness

	/*
		This additon makes it run slower.
		if result.Imag {
			d.Fitness = fitness * 1.1
		} else {
			d.Fitness = fitness * 0.9
		}
	*/
}

// Residual Sum of Squares
func calcDifference(gen []float64, target []float64) (float64, error) {
	if len(gen) == 0 || len(target) == 0 {
		fmt.Printf("error in generation of summarize")
		return 9999.0, ErrNullSummarize
	}
	var d float64
	for i, v := range padSlice(gen, target) {
		d += squareDifference(v, target[i])
	}

	if math.IsNaN(math.Sqrt(d)) {
		fmt.Printf("NAN detected.\n")
		return 9999.0, ErrNaNFitness
	}

	return math.Sqrt(d), nil
}

func squareDifference(x, y float64) float64 {
	d := x - y
	return d * d
}

// Sometimes the values returned from summarize aren't the correct length since
// there is a cutoff for those close to 0. Pad each slice to be the same length.
func padSlice(s []float64, target []float64) []float64 {
	if len(s) == len(target) {
		return s
	} else {
		for i := len(s); i < len(target); i++ {
			s = append(s, 0)
		}
	}

	return s
}

func createPopulation() (population []Organism) {
	var wg sync.WaitGroup
	population = make([]Organism, *PopSize)

	sema := make(chan struct{}, 4)

	for i := 0; i < *PopSize; i++ {
		sema <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() {
				<-sema
				wg.Done()
			}()
			org := CreateOrganism(*NumAtoms)
			population[i] = org
		}(i)
	}
	wg.Wait()

	return
}

func createPool(population []Organism, target []float64) (pool []Organism) {
	pool = make([]Organism, 0)
	// get top best fitting organisms
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	// fmt.Printf("The top fitness of the pool is: %v %v %v", population[0].Fitness, population[1].Fitness, population[2].Fitness)

	// This is what fraction survives to the next generation.
	//	fmt.Println("length population, ", len(population))
	fraction := *PoolSize * float64(*PopSize)

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

func delFolders(o []Organism, topOrganism Organism) {
	for _, v := range o {
		if v.Path == topOrganism.Path {
			continue
		} else {
			os.RemoveAll(path.Dir(v.Path))
		}
	}
}

// perform natural selection to create the next generation
// We should use a weighted method.
func naturalSelection(pool []Organism, population []Organism, target []float64) []Organism {
	// Children = [pool + empty slice]; len = population.
	next := make([]Organism, len(population)-len(pool))

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

		a, b := tournamentRound(pool)

		child := crossover(a, b)
		child.mutate()

		err := child.saveToFile(*NumAtoms)
		if err != nil {
			panic(err)
		}

		child.calcFitness()

		children[i] = child
	}

	//fmt.Println("The length of next is: ", len(children))
	return children
}

func tournamentRound(pool []Organism) (d Organism, m Organism) {
	// This grabs three random organisms from the pool and sorts them.
	makeBracket := func() []Organism {
		round := make([]Organism, 0)

		for i := 0; i < *TournamentPool; i++ {
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
	a, b := makeBracket(), makeBracket()

	// Since these slices are sorted, the first element is the most fit organism.
	return a[0], b[0]
}

// crosses over 2 Organism strings
// We should use the blending method prescribed in Haupt and Haupt.
// Furthermore, each cross should produce two offspring.
func crossover(d1 Organism, d2 Organism) Organism {
	childDNA := make(DNA, len(d1.DNA))

	for i, pChr := range d1.DNA {
		ch := make(Chromosome, len(pChr))
		childDNA[i] = ch
	}

	child := Organism{
		DNA:     childDNA,
		Path:    "",
		Fitness: 0,
	}

	// if rand.Float64() <= *CrossOverRate {
	// Points to the left come from the first parent.
	// Points to the right come from the other parent.
	// Points in the middle are blended.

	for i, chr := range d1.DNA {
		mid := rand.Intn(len(chr))
		for j := 0; j < len(chr); j++ {
			if j < mid {
				child.DNA[i][j] = d1.DNA[i][j]
			} else if i > mid {
				child.DNA[i][j] = d2.DNA[i][j]
			} else if i == mid {
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

// Mutation function is unclear. I had it previously generate a new random number, but now it'll add or subtract.
func (o *Organism) mutate() {
	for c, chr := range o.DNA {
		for i := 0; i < len(chr); i++ {
			chance := rand.Float64()
			if chance <= *MutationRate {
				o.DNA[c][i] = rand.Float64()
				if RandBool() {
					o.DNA[c][i] = -o.DNA[c][i]
				}

				if chance <= *ZeroChance {
					o.DNA[c][i] = 0.0
				}
			}
		}
	}
}

func getBest(population []Organism) Organism {
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	return population[0]
}

var OutputPath string

func init() {
	outputPath, err := filepath.Abs(*OutFile)
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
}

func main() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	population := createPopulation()

	found := false
	generation := 0
	for !found {
		generation++
		bestOrganism := getBest(population)
		if bestOrganism.Fitness < *FitnessLimit {
			found = true

			f, err := os.OpenFile(OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			foundString := fmt.Sprintf("The path to the best organism is %v\n", bestOrganism.Path)

			if _, err = f.WriteString(foundString); err != nil {
				panic(err)
			}

			f.Close()

			bestPath := "best/final"
			bestErr := bestOrganism.SaveBestOrganism(*NumAtoms, bestPath)
			if bestErr != nil {
				fmt.Printf("Erro saving best organism, %v\n", err)
			}

			elapsed := time.Since(start)
			fmt.Printf("\nTotal time taken: %s\n", elapsed)

			return

		} else {
			pool := createPool(population, TargetFrequencies)

			population = naturalSelection(pool, population, TargetFrequencies)

			if generation%10 == 0 {
				sofar := time.Since(start)

				summaryStep := fmt.Sprintf("The path to the best organism is %v.\n \nTime taken so far: %s | generation: %d | fitness: %f | pool size: %d\n", bestOrganism.Path, sofar, generation, bestOrganism.Fitness, len(pool))

				f, err := os.OpenFile(OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}

				if _, err = f.WriteString(summaryStep); err != nil {
					panic(err)
				}

				f.Close()

				bestPath := fmt.Sprintf("best/%d", generation)
				bestErr := bestOrganism.SaveBestOrganism(*NumAtoms, bestPath)
				if bestErr != nil {
					fmt.Printf("Erro saving best organism, %v\n", err)
				}

				if generation >= *GenLimit {
					f, err := os.OpenFile(OutputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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
