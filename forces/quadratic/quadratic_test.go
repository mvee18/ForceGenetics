package quadratic

import (
	"fmt"
	"testing"
)

func TestLinearInterpolation(t *testing.T) {
	t.Run("testing linear interpolation recursion with disparate parents", func(t *testing.T) {
		iter := 0.0

		lin, err := LinearInterpolation(&iter, 0.67, 0.31, 0.62)
		if err != nil {
			t.Error(err)
		}

		fmt.Println(lin)

	})

	t.Run("testing too large beta", func(t *testing.T) {

	})
}
