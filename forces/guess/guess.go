package guess

import (
	"bufio"
	"ga/forces/models"
	"ga/forces/utils"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	s1 := rand.NewSource(time.Now().UnixNano() + 2561)
	r1 = rand.New(s1)
}

var mu sync.Mutex
var r1 *rand.Rand

func ReadDNA(dirPath string) models.DNA {
	fort15 := path.Join(dirPath, "fort.15")
	fort30 := path.Join(dirPath, "fort.30")
	fort40 := path.Join(dirPath, "fort.40")

	DNA := make(models.DNA, 3)

	DNA[0] = readFortFile(fort15)
	DNA[1] = readFortFile(fort30)
	DNA[2] = readFortFile(fort40)

	return DNA
}

func readFortFile(fp string) []float64 {
	absPath, err := filepath.Abs(fp)
	if err != nil {
		log.Fatalf("Could not open input file, %v\n", err)
	}

	f, err := os.Open(absPath)
	if err != nil {
		log.Fatalf("Could not open input file, %v\n", err)
	}

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example

	chromosome := make([]float64, 0)

	count := 0
	for scanner.Scan() {
		// Skip the header.
		if count == 0 {
			count++
		} else {
			fcs := strings.Fields(scanner.Text())
			fcsFloat := parseText(fcs)

			for _, fc := range fcsFloat {
				chromosome = append(chromosome, fc)
			}

			// fmt.Println("Parsed Line: ", fcsFloat)

		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return chromosome
}

func parseText(s []string) []float64 {
	var freq []float64

	for _, v := range s {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}

		freq = append(freq, value)
	}

	return freq
}

func MockB3LYP(dirPath string) models.DNA {
	DNA := ReadDNA(dirPath)

	newDNA := mutateDNA(DNA)

	// newOrg := models.Organism{
	// 	DNA:     newDNA,
	// 	Fitness: 0.0,
	// 	Path:    "",
	// }

	// newOrg.SaveToFile(*flags.NumAtoms)
	// newOrg.CalcFitness()

	// fmt.Printf("the fitness of the new org is %v\n", newOrg.Fitness)

	return newDNA
}

// This function will randomly mutate the DNA by +- 0.5 to simulate the
// lower level calculations.
func mutateDNA(d models.DNA) models.DNA {

	var wg sync.WaitGroup
	sema := make(chan struct{}, 4)

	newDNA := make(models.DNA, len(d))
	newDNA[0] = make(models.Chromosome, len(d[0]))
	newDNA[1] = make(models.Chromosome, len(d[1]))
	newDNA[2] = make(models.Chromosome, len(d[2]))

	for i, v := range d {
		for j, gene := range v {
			sema <- struct{}{}
			wg.Add(1)
			go func(i, j int, gene float64) {
				defer func() {
					<-sema
					wg.Done()
				}()

				if utils.RandBool() {
					newVal := gene + (0.0 + rand.Float64()*(0.1-0.0))
					incrementAndCheck(&newVal, i+2)
					newDNA[i][j] = newVal

				} else {
					newVal := gene - (0.0 + rand.Float64()*(0.1-0.0))
					incrementAndCheck(&newVal, i+2)
					newDNA[i][j] = newVal
				}

			}(i, j, gene)
		}
	}

	return newDNA
}

func incrementAndCheck(v *float64, dn int) {
	dom := utils.SelectDomain(dn)

	// If the value is out of bounds...
	for math.Abs(*v) > dom {
		if *v > 0 {
			// If positive, subtract until it is below...
			mu.Lock()
			*v -= (0.0 + rand.Float64()*(0.1-0.0))
			mu.Unlock()

		} else {
			// Else, if negative, add until within the bounds.
			mu.Lock()
			*v += (0.0 + rand.Float64()*(0.1-0.0))
			mu.Unlock()
		}
	}

}
