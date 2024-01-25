package simplex_test

import (
	"testing"

	"github.com/CGuilhem/Perlin-Noise/Go/simplex/simplex"
)

func BenchmarkNoise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simplex := simplex.NewSimplex()
		simplex.Noise(0.1, 0.1)
	}
}
