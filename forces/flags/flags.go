package flags

import (
	"flag"
	"fmt"
	"ga/forces/quotes"
	"math/rand"
	"testing"
	"time"
)

var (
	MutationRate         = flag.Float64("mut", 0.10, "mutation rate as a decimal for TGA core")
	MutationRateInformed = flag.Float64("imut", 0.02, "mutation rate as a decimal for informed GA core")
	PopSize              = flag.Int("pop", 600, "population size")
	NumAtoms             = flag.Int("n", 4, "number of atoms")
	PoolSize             = flag.Float64("pool", 0.50, "fraction size of the the previous generation that survives")
	FitnessLimit         = flag.Float64("f", 1.0, "fitness criteria")
	OutFile              = flag.String("o", "forces.out", "name of output file")
	PathToSpectro        = flag.String("sp", "./spectro", "path/to/spectro")
	PathToSpectroIn      = flag.String("i", "./spectro.in", "path/to/spectro.in")
	ZeroChance           = flag.Float64("z", 0.02, "chance mutation will set value to zero instead of adding/subtracting")
	TournamentPool       = flag.Int("t", 150, "the number of organisms selected to compete to be parents")
	DerivativeLevel      = flag.Int("d", 4, "this is the level of derivative")
	// TODO: This path doesn't work. It defaults to the directory that the script is being executed in.
	Fort15File          = flag.String("ft2", "./fort.15", "path to the fort.15 file that will be used for 3rd derivatives")
	GenLimit            = flag.Int("l", 100000, "the maximum number of generations")
	PseudoCrossOverRate = flag.Float64("pc", 0.50, "the chance that two parents chromosomes will crossover in the psuedo GA.")
	CrossOverRate       = flag.Float64("c", 0.95, "the chance that two parents chromosomes will crossover, otherwise the only possible change would be mutations.")
	FreqInputFile       = flag.String("fi", "forces.inp", "the input from where the frequencies will be read")
	Domain15            = flag.Float64("dom15", 1.0, "the starting guess for the maximum values for the cartesian force constants")
	Domain30            = flag.Float64("dom30", 3.0, "the starting guess for the maximum values for the cartesian force constants")
	Domain40            = flag.Float64("dom40", 10.0, "the starting guess for the maximum values for the cartesian force constants")
)

func init() {
	var _ = func() bool {
		testing.Init()
		return true
	}()

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("%s\n", quotes.Sayings[rand.Intn(len(quotes.Sayings))])

	fmt.Printf("The pop size is %d, the pool size is %f\n", *PopSize, *PoolSize)
}
