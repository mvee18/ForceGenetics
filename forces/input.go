package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
