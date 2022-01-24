package informed

import (
	"bytes"
	"fmt"
	"ga/forces/flags"
	"ga/forces/models"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ntBre/chemutils/summarize"
)

func TestSummarize(t *testing.T) {
	d := models.Organism{
		DNA:     models.DNA{},
		Fitness: 0.0,
		Path:    "testorg/fort.40",
	}

	t.Run("testing quartile organism generation", func(t *testing.T) {
		f, err := os.Open(d.Path)
		defer f.Close()

		input, err := os.Open(*flags.PathToSpectroIn)
		if err != nil {
			log.Fatalf("Error opening file of spectro in, %v\n", err)
		}

		// fmt.Println(string(b))

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

		r := bytes.NewReader(outBytes)
		result := summarize.Spectro(r)

		// fmt.Printf("%#v", result)
		// fmt.Println(d.Path)

		if !reflect.DeepEqual(result.LX, []float64{5764.6, 4724.45, 2917.15, 2484.36, 1519.33, 5320.25, 148.77, 770.87, 1181.23, 2140.61, 4340.03, 1019.58}) {
			t.Errorf("wrong LX matrix")
		}
	})
}

func TestGenerateOrganism(t *testing.T) {
	org := CreateInformedOrganism(4, true)

	fmt.Println(org.Fitness)

	os.RemoveAll(path.Dir(org.Path))
}

func TestDirectedMutation(t *testing.T) {

	i := InformedOrganism{
		Organism: models.Organism{
			DNA: models.DNA{
				{1, 2, 3},
				{4, 5, 6, 7},
				{8, 9, 10, 11, 12},
			},

			Fitness: 0.0,
			Path:    "",
		},

		Direction: [3][]bool{
			{true, true, true},
			{true, true, false, false},
			{false, true, false, true, false},
		}}

	t.Run("test directed mutation with mocked function", func(t *testing.T) {
		f := DirectedMutation(i, (InformedOrganism).MockCalcFitness)
		// DirectedMutation(&i, (*InformedOrganism).CalcFitness)

		want := [3][]bool{
			{false, false, false},
			{false, false, true, true},
			{true, false, true, false, true},
		}

		if !reflect.DeepEqual(f.Direction, want) {
			t.Errorf("error in changing direction, got %v\n, wanted %v\n", i.Direction, want)
		}

		// fmt.Println(i.DNA)
	})
}

// Always return worse fitness.
func (i InformedOrganism) MockCalcFitness() float64 {
	return 10000
}

func TestIncrementAndCheck(t *testing.T) {
	t.Run("test with positive number", func(t *testing.T) {
		v := 50.0
		dn := 4

		incrementAndCheck(&v, dn)

		wantedLimit := *flags.Domain40
		if v > wantedLimit {
			t.Errorf("wanted less than %v, got %v\n", wantedLimit, v)
		}

		fmt.Printf("new v: %v\n", v)
	})

	t.Run("test with negative number", func(t *testing.T) {
		v := -50.0
		dn := 4

		incrementAndCheck(&v, dn)

		wantedLimit := *flags.Domain40
		if v < -wantedLimit {
			t.Errorf("wanted more than %v, got %v\n", wantedLimit, v)
		}

		fmt.Printf("new v: %v\n", v)
	})
}

func TestRunIGA(t *testing.T) {
	t.Run("run iga", func(t *testing.T) {
		imm := make(chan models.OrganismAndBias)

		RunInformedGA(imm)
	})
}
