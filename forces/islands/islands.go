package islands

import (
	"fmt"
	"ga/forces/informed"
	"ga/forces/models"
	"ga/forces/pseudo"
	trad "ga/forces/traditional"
)

// Not exactly sure about the methods they share in common.
// Might want to add a TPopulation for the TGA.
type Island interface {
}

func RunIslands(imm chan models.Organism, mig chan models.Organism) {
	go trad.RunTGA(imm, mig)

	go informed.RunInformedGA(imm, mig)

	go pseudo.RunPGA(imm, mig)

	for i := 0; i < 3; i++ {
		select {
		case <-mig:
			fmt.Println(i)
		}
	}
}
