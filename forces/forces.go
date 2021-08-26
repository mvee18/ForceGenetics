package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
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

var TargetFrequencies = []float64{
	3943, 3832, 1649,
	0, 0, 0,
	0, 0, 0,
}

// The organism is going to be the array of force constants.
// We should be able to represent this as a one dimensions array,
// then reconstruct it to run it in spectro.
type Organism struct {
	DNA     []float64
	Path    string
	Fitness float64
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func GetNumForceConstants(n int, dn int) int {
	if dn == 2 {
		return int(math.Pow(float64(n), 2)) * 3 * 3
	} else {
		return 0
	}
}

func CreateOrganism(numAtoms int) (organism Organism) {
	organismSize := GetNumForceConstants(numAtoms, 2)
	ba := make([]float64, organismSize)
	for i := 0; i < organismSize; i++ {
		ba[i] = (rand.Float64())
		if RandBool() {
			ba[i] = -ba[i]
		}
	}

	organism = Organism{
		DNA:     ba,
		Path:    "",
		Fitness: 0,
	}

	err := organism.saveToFile(numAtoms)
	if err != nil {
		fmt.Printf("Error in saving organism to file, %v\n", err)
	}

	organism.calcFitness(TargetFrequencies)

	return organism
}

// Before we can calc the fitness, we have to save the files so that spectro can use them.
// Let's use temp files.
func (d *Organism) saveToFile(natoms int) error {
	dir, err := ioutil.TempDir(".", "forceOrganisms")
	if err != nil {
		log.Fatalf("Could not open temp dir, %v\n", err)
		return err
	}

	tempfn := path.Join(dir, "fort.15")
	organismFile, err := os.Create(tempfn)
	if err != nil {
		log.Fatalf("Could not open temp file, %v\n", err)
		return err
	}

	// Now we need to format the file correctly.
	// Spectro is 20.12f
	fmt.Fprintf(organismFile, "%5d%5d", natoms, 2*natoms*natoms)
	for i := range d.DNA {
		if i%3 == 0 {
			fmt.Fprintf(organismFile, "\n")
		}
		fmt.Fprintf(organismFile, "%20.10f", d.DNA[i])
	}
	organismFile.Write([]byte("\n"))

	d.Path = organismFile.Name()

	organismFile.Close()

	return nil
}

func (d *Organism) SaveBestOrganism(natoms int, filePath string) error {
	err := os.MkdirAll(filePath, 0700)
	if err != nil {
		log.Fatalf("Could not open temp dir, %v\n", err)
		return err
	}

	tempfn := path.Join(filePath, "fort.15")
	organismFile, err := os.Create(tempfn)
	if err != nil {
		log.Fatalf("Could not open temp file, %v\n", err)
		return err
	}

	// Now we need to format the file correctly.
	// Spectro is 20.12f
	fmt.Fprintf(organismFile, "%5d%5d", natoms, 6*natoms)
	for i := range d.DNA {
		if i%3 == 0 {
			fmt.Fprintf(organismFile, "\n")
		}
		fmt.Fprintf(organismFile, "%20.10f", d.DNA[i])
	}
	organismFile.Write([]byte("\n"))

	d.Path = organismFile.Name()

	organismFile.Close()

	return nil
}

// To calculate the fitness, we must run it through spectro.
// We can save the results to a temp file and get a difference.
// The smaller the differences, the greater the fitness (1/difference).
func (d *Organism) calcFitness(target []float64) {
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

	outString := string(outBytes)
	temp := strings.Split(outString, "\n")

	flag := false
	lxmflag := false
	var lxm []string
	for _, v := range temp {
		if lxmflag {
			if flag {
				if strings.Contains(v, "---") {
					flag = false
				} else {
					lxm = append(lxm, v)
				}

			} else if strings.Contains(v, "") && !flag {
				if strings.Contains(v, "---") {
					lxmflag = false
				} else {
					lxm = append(lxm, v)
				}
			}
		}

		if strings.Contains(v, "LXM") {
			flag = true
			lxmflag = true
			lxm = append(lxm, v)
		}

	}

	endBytes := func(list []string) int {
		return len(list) - 4
	}

	firstLine := lxm[1:4]

	secondLine := lxm[endBytes(lxm):]

	fields := func(s []string) []string {
		var new []string
		for _, v := range s {
			new = strings.Fields(v)
		}

		return new
	}

	trimFirst := fields(firstLine)
	trimSecond := fields(secondLine)

	var newLXM []string

	if !reflect.DeepEqual(trimFirst, trimSecond) {
		newLXM = append(trimFirst, trimSecond...)
	} else {
		newLXM = trimFirst
	}

	var LXMfloat []float64

	for _, v := range newLXM {
		stringToFloat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatalf("error parsing float, %v\n", err)
		}

		LXMfloat = append(LXMfloat, stringToFloat)
	}

	fitness := calcDifference(LXMfloat)
	if fitness == 0 {
		d.Fitness = 1
	} else {
		d.Fitness = fitness
	}
}

func calcDifference(lxm []float64) float64 {
	var d float64
	for i, v := range lxm {
		d += squareDifference(v, TargetFrequencies[i])
	}

	return math.Sqrt(d)
}

func squareDifference(x, y float64) float64 {
	d := x - y
	return d * d
}

func createPopulation() (population []Organism) {
	population = make([]Organism, *PopSize)
	for i := 0; i < *PopSize; i++ {
		population[i] = CreateOrganism(*NumAtoms)
	}
	return
}

func createPool(population []Organism, target []float64) (pool []Organism) {
	pool = make([]Organism, 0)
	// get top best fitting organisms
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})
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

	for i := len(pool); i < len(population); i++ {
		/*
			r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
			a := pool[r1]
			b := pool[r2]
		*/

		a, b := tournamentRound(pool)

		child := crossover(a, b)
		child.mutate()

		child.saveToFile(*NumAtoms)

		child.calcFitness(target)

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
			sort.SliceStable(round, func(i, j int) bool {
				return round[i].Fitness < round[j].Fitness
			})
		}

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
	childDNA := make([]float64, len(d1.DNA))
	child := Organism{
		DNA:     childDNA,
		Path:    "",
		Fitness: 0,
	}

	// Points to the left come from the first parent.
	// Points to the right come from the other parent.
	// Points in the middle are blended.

	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i < mid {
			child.DNA[i] = d1.DNA[i]
		} else if i > mid {
			child.DNA[i] = d2.DNA[i]
		} else if i == mid {
			child.DNA[i] = crossOverA(d1.DNA[i], d2.DNA[i])
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
	for i := 0; i < len(o.DNA); i++ {
		chance := rand.Float64()
		if chance <= *MutationRate {
			o.DNA[i] = rand.Float64()
			if RandBool() {
				o.DNA[i] = -o.DNA[i]
			}

			if chance <= *ZeroChance {
				o.DNA[i] = 0.0
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
			}

			delFolders(pool, bestOrganism)
			delFolders(population, bestOrganism)
		}

	}

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
