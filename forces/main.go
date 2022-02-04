package main

import (
	"ga/forces/islands"
	"ga/forces/models"
)

var (
	immPGA = make(chan models.Migrant)
	immTGA = make(chan models.Migrant)
	immIGA = make(chan models.Migrant)
)

func main() {
	islands.RunIslands(immPGA, immTGA, immIGA)
}
