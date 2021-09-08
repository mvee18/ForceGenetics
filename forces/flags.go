package main

import (
	"flag"
	"fmt"
	"testing"
)

var (
	MutationRate    = flag.Float64("mut", 0.04, "mutation rate as a decimal")
	PopSize         = flag.Int("pop", 600, "population size")
	NumAtoms        = flag.Int("n", 6, "number of atoms")
	PoolSize        = flag.Float64("pool", 0.50, "fraction size of the the previous generation that survives")
	FitnessLimit    = flag.Float64("f", 1.0, "fitness criteria")
	OutFile         = flag.String("o", "forces.out", "name of output file")
	PathToSpectro   = flag.String("sp", "./spectro", "path/to/spectro")
	PathToSpectroIn = flag.String("i", "./spectro.in", "path/to/spectro.in")
	ZeroChance      = flag.Float64("z", 0.02, "chance mutation will set value to zero instead of adding/subtracting")
	TournamentPool  = flag.Int("t", 3, "the number of chromosomes selected to compete to be parents")
	DerivativeLevel = flag.Int("d", 2, "this is the level of derivative")
	Fort15File      = flag.String("ft2", "./fort.15", "path to the fort.15 file that will be used for 3rd derivatives")
)

func init() {
	var _ = func() bool {
		testing.Init()
		return true
	}()

	flag.Parse()
	fmt.Printf("The pop size is %d, the pool size is %f\n", *PopSize, *PoolSize)
}
