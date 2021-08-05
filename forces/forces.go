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
	"sort"
	"strconv"
	"strings"
	"time"
)

// MutationRate is the rate of mutation
var MutationRate = 0.0004

// PopSize is the size of the population (was 500)
var PopSize = 300

// The number of atoms is critical to the size of the organism.
var NumAtoms = 6

// Size of breeding pool (was 30)
var PoolSize = 35

var FitnessLimit = 1000.0

var TargetFrequencies = []float64{
	820.24, 804.08, 737.75,
	580.87, 573.06, 525.60,
	363.23, 274.64, 194.26,
	170.33, 170.26, 0.34,
	0.00, 0.00, 0.00,
	0.00, 0.00, 0.00,
}

var PathToSpectro = "./spectro"

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
	fmt.Fprintf(organismFile, "%5d%5d", natoms, natoms*natoms)
	for i := range d.DNA {
		if i%3 == 0 {
			fmt.Fprintf(organismFile, "\n")
		}
		fmt.Fprintf(organismFile, "%20.10f", d.DNA[i])
	}
	organismFile.Write([]byte("\n"))

	d.Path = organismFile.Name()

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

	input, err := os.Open("./spectro.in")
	if err != nil {
		log.Fatalf("Error opening file of spectro in, %v\n", err)
	}

	spectroAbs, err := filepath.Abs(PathToSpectro)
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

	newLXM := append(trimFirst, trimSecond...)

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
	}

	d.Fitness = fitness
}

func calcDifference(lxm []float64) float64 {
	var d float64
	for i, v := range lxm {
		d += squareDifference(v, TargetFrequencies[i])
	}

	return math.Sqrt(float64(d))
}

func squareDifference(x, y float64) float64 {
	d := x - y
	return d * d
}

func createPopulation() (population []Organism) {
	population = make([]Organism, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = CreateOrganism(NumAtoms)
	}
	return
}

func createPool(population []Organism, target []float64) (pool []Organism) {
	pool = make([]Organism, 0)
	// get top best fitting organisms
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})
	top := population[0 : PoolSize+1]
	bottom := population[PoolSize+2:]

	delBottomFolders(bottom)

	// if there is no difference between the top  organisms, the population is stable
	// and we can't get generate a proper breeding pool so we make the pool equal to the
	// population and reproduce the next generation
	if top[len(top)-1].Fitness-top[0].Fitness == 0 {
		pool = population
		return
	}
	// create a pool for next generation
	for i := 0; i < len(top)-1; i++ {
		num := (top[PoolSize].Fitness - top[i].Fitness)
		for n := 0.0; n < num; n++ {
			pool = append(pool, top[i])
		}
	}
	return
}

func delBottomFolders(o []Organism) {
	for _, v := range o {
		defer os.RemoveAll(path.Dir(v.Path))
	}
}

// perform natural selection to create the next generation
func naturalSelection(pool []Organism, population []Organism, target []float64) []Organism {
	next := make([]Organism, len(population))

	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		child := crossover(a, b)
		child.mutate()

		child.saveToFile(NumAtoms)

		child.calcFitness(target)

		next[i] = child
	}
	return next
}

// crosses over 2 Organism strings
func crossover(d1 Organism, d2 Organism) Organism {
	childDNA := make([]float64, len(d1.DNA))
	child := Organism{
		DNA:     childDNA,
		Path:    "",
		Fitness: 0,
	}

	mid := rand.Intn(len(d1.DNA))
	for i := 0; i < len(d1.DNA); i++ {
		if i > mid {
			child.DNA[i] = d1.DNA[i]
		} else {
			child.DNA[i] = d2.DNA[i]
		}

	}

	return child
}

func (o *Organism) mutate() {
	for i := 0; i < len(o.DNA); i++ {
		if rand.Float64() < MutationRate {
			o.DNA[i] = (rand.Float64())
			if RandBool() {
				o.DNA[i] = -o.DNA[i]
			}
		}
	}
}

func getBest(population []Organism) Organism {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
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
		if bestOrganism.Fitness < FitnessLimit {
			found = true
			fmt.Printf("The path to the best organism is %v\n", bestOrganism.Path)
		} else {
			pool := createPool(population, TargetFrequencies)
			population = naturalSelection(pool, population, TargetFrequencies)
			if generation%10 == 0 {
				sofar := time.Since(start)
				fmt.Printf("The path to the best organism is %v\n", bestOrganism.Path)
				fmt.Printf("\nTime taken so far: %s | generation: %d | fitness: %f | pool size: %d", sofar, generation, bestOrganism.Fitness, len(pool))
			}
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("\nTotal time taken: %s\n", elapsed)
}
