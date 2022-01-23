package main

import (
	"ga/forces/islands"
	"ga/forces/models"
)

var (
	immPGA = make(chan models.Organism)
	immTGA = make(chan models.Organism)
	immIGA = make(chan models.Organism)
)

func main() {
	islands.RunIslands(immPGA, immTGA, immIGA)
}
