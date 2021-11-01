package main

import (
	"fmt"
	"ga/forces/selection"
	"testing"
)

func TestGenerateOrganism(t *testing.T) {
	t.Run("testing generate organism", func(t *testing.T) {
		org := selection.CreateOrganism(4)

		fmt.Println(org.Fitness)
	})
}
