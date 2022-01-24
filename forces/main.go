package main

import (
	"ga/forces/islands"
	"ga/forces/models"
)

var (
	immPGA = make(chan models.OrganismAndBias)
	immTGA = make(chan models.OrganismAndBias)
	immIGA = make(chan models.OrganismAndBias)
)

func main() {
	islands.RunIslands(immPGA, immTGA, immIGA)
}
