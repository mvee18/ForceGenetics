package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/ntBre/chemutils/summarize"
)

func TestSummarize(t *testing.T) {
	f, err := os.Open("spectro.out")
	if err != nil {
		t.Fatalf("failed to open file")
	}

	result := summarize.Spectro(f)

	fmt.Printf("%#v", result)
}
