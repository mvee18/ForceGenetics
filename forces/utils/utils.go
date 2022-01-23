package utils

import (
	"bufio"
	"errors"
	"fmt"
	"ga/forces/flags"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var ErrNaNFitness = errors.New("not a number fitness")
var ErrNullSummarize = errors.New("null value of organism")

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

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func CalcDifference(gen []float64, target []float64) (float64, error) {
	if len(gen) == 0 || len(target) == 0 {
		fmt.Printf("gen: %v\n target: %v\n", gen, target)
		fmt.Printf("error in generation of summarize")
		return 9999.0, ErrNullSummarize
	}
	var d float64
	for i, v := range padSlice(gen, target) {
		d += squareDifference(v, target[i])
	}

	if math.IsNaN(math.Sqrt(d)) {
		// fmt.Printf("NAN detected.\n")
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

func ReadInput(fp string) (harmonic []float64, rot []float64, fund []float64) {
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

	harmonic, rot, fund = make([]float64, 0), make([]float64, 0), make([]float64, 0)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "HARMONICS:") {
			harmSub := strings.Fields(scanner.Text())[1:]
			harmonic = parseText(harmSub)

			fmt.Println("HARMONICS: ", harmonic)

		} else if strings.Contains(scanner.Text(), "ROTATIONAL:") {
			rotSub := strings.Fields(scanner.Text())[1:]
			rot = parseText(rotSub)

			fmt.Println("ROTS: ", rotSub)

		} else if strings.Contains(scanner.Text(), "FUND:") {
			fundSub := strings.Fields(scanner.Text())[1:]
			fund = parseText(fundSub)

			fmt.Println("FUNDS: ", fundSub)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return harmonic, rot, fund
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

func RandValueDomain(dn int) float64 {
	switch dn {
	case 2:
		return (0.0 + rand.Float64()*(*flags.Domain15-0.0))
	case 3:
		return (0.0 + rand.Float64()*(*flags.Domain30-0.0))
	case 4:
		return (0.0 + rand.Float64()*(*flags.Domain40-0.0))
	default:
		panic("undefined derivative level: could not select domain.")
	}
}

func SelectDomain(dn int) float64 {
	switch dn {
	case 2:
		return *flags.Domain15
	case 3:
		return *flags.Domain30
	case 4:
		return *flags.Domain40
	default:
		panic("undefined derivative level: could not select domain.")
	}
}

func NewOutputFile(path string) string {
	rootDir := filepath.Dir(*flags.OutFile)
	newDir := filepath.Join(rootDir, path)

	return newDir
}
