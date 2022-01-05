package models

import (
	"bytes"
	"errors"
	"fmt"
	"ga/forces/flags"
	"ga/forces/utils"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/ntBre/chemutils/summarize"
)

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

// This interface allow us to use the same methods for the informed
// and the TGA?
type Organismer interface {
	SaveToFile(natoms int) error
	CalcFitness()
}

var FortFiles = []string{"fort.15", "fort.30", "fort.40"}

var TargetFrequencies = []float64{}

var TargetRotational = []float64{}

var TargetFund = []float64{}

var ErrCalcFitness = errors.New("error calculating fitness of organism")
var ErrCalcHarmFitness = errors.New("error calculating harm fitness")
var ErrCalcRotFitness = errors.New("error calculating rot fit")
var ErrCalcFundFitness = errors.New("error calculating fund fit")

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

		fmt.Fprintf(organismFile, "%5d%5d", natoms, utils.GetNumForceConstants(natoms, i+2))
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
func (d *Organism) CalcFitness() {
	// Let's begin with opening the fort file in each organism.
	f, err := os.Open(d.Path)
	if err != nil {
		log.Fatalf("Error opening file of path, %v\n", err)
	}

	defer f.Close()

	input, err := os.Open(*flags.PathToSpectroIn)
	if err != nil {
		log.Fatalf("Error opening file of spectro in, %v\n", err)
	}

	defer input.Close()

	spectroAbs, err := filepath.Abs(*flags.PathToSpectro)
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

	ParseOutput(d, outBytes, *flags.DerivativeLevel)
}

// Mutation function is unclear. There are many ideas in the gaussian crossover paper.
func (o *Organism) Mutate() {
	for c, chr := range o.DNA {
		for i := 0; i < len(chr); i++ {
			chance := rand.Float64()
			if chance <= *flags.MutationRate {
				o.DNA[c][i] = utils.RandValueDomain(c + 2)
				if utils.RandBool() {
					o.DNA[c][i] = -o.DNA[c][i]
				}

				if chance <= *flags.ZeroChance {
					o.DNA[c][i] = 0.0
				}
			}
		}
	}
}

// Before we can calc the fitness, we have to save the files so that spectro can use them.
// Let's use temp files.
func (d *Organism) SaveToFile(natoms int) error {
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

		fmt.Fprintf(organismFile, "%5d%5d", natoms, utils.GetNumForceConstants(natoms, i+2))
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

// TODO: Add a constraint on negative frequencies. Should reduce the fitness of
// the organism, even if the freq get closer.
func ParseOutput(d *Organism, by []byte, derivative int) {
	r := bytes.NewReader(by)
	result := summarize.Spectro(r)

	// fmt.Printf("%#v", result)
	// fmt.Println(d.Path)

	fitness := 9999.0

	var err error

	switch derivative {
	case 2:
		fitness, err = utils.CalcDifference(result.LX, TargetFrequencies)
		if err != nil {
			if err == utils.ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Fatalln(ErrCalcHarmFitness)
			}
		}
	case 3:
		harmFitness, err := utils.CalcDifference(result.Harm, TargetFrequencies)
		if err != nil {
			if err == utils.ErrNullSummarize {
				fitness = 9999.99
			} else {
				log.Fatalln(ErrCalcHarmFitness)
			}
		}

		rotFitness, err := utils.CalcDifference(result.Rots[0], TargetRotational)
		if err != nil {
			if err == utils.ErrNullSummarize {
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

		harmFitness, err := utils.CalcDifference(result.Harm, TargetFrequencies)
		if err != nil {
			if err == utils.ErrNullSummarize {
				fitness = 9999.99
			} else {
				// log.Printf("%v in organism %v\n", err, d.Path)
				d.Fitness = 99999.99
				return
			}
		}

		rotFitness, err := utils.CalcDifference(result.Rots[0], TargetRotational)
		if err != nil {
			if err == utils.ErrNullSummarize {
				fitness = 9999.99
			} else {
				// log.Printf("%v in organism %v\n", err, d.Path)
				d.Fitness = 99999.99
				return
			}
		}

		fundFitness, err := utils.CalcDifference(result.Fund, TargetFund)
		if err != nil {
			if err == utils.ErrNullSummarize {
				fitness = 9999.99
			} else {
				// log.Printf("%v in organism %v\n", err, d.Path)
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
