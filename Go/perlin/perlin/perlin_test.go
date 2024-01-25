package perlin_test

import (
	"testing"

	"github.com/CGuilhem/Perlin-Noise/Go/perlin/perlin"
)

func BenchmarkPerlin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		perlin.Perlin(0.0, 0.0, nil)
	}
}
