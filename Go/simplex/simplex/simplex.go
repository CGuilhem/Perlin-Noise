package simplex

import (
	"math"
	"math/rand"
	"time"
)

type Simplex struct {
	grad3 [12][3]int
	p     [256]int
	perm  [512]int
}

func (s *Simplex) dot(g [3]int, x, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}

func (s *Simplex) Noise(xin, yin float64) float64 {
	F2 := 0.5 * (math.Sqrt(3.0) - 1.0)
	l := (xin + yin) * F2
	i := math.Floor(xin + l)
	j := math.Floor(yin + l)

	G2 := (3.0 - math.Sqrt(3.0)) / 6.0
	t := (i + j) * G2
	X0 := i - t
	Y0 := j - t
	x0 := xin - X0
	y0 := yin - Y0

	var i1, j1 int
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2
	y2 := y0 - 1.0 + 2.0*G2

	ii := int(i) & 255
	jj := int(j) & 255
	gi0 := s.perm[ii+s.perm[jj]] % 12
	gi1 := s.perm[ii+i1+s.perm[jj+j1]] % 12
	gi2 := s.perm[ii+1+s.perm[jj+1]] % 12

	t0 := 0.5 - x0*x0 - y0*y0
	var n0 float64
	if t0 < 0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * s.dot(s.grad3[gi0], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	var n1 float64
	if t1 < 0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * s.dot(s.grad3[gi1], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	var n2 float64
	if t2 < 0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * s.dot(s.grad3[gi2], x2, y2)
	}

	return 70.0 * (n0 + n1 + n2)
}

func NewSimplex() *Simplex {
	rand.Seed(time.Now().UnixNano())

	simplex := &Simplex{}

	for i := 0; i < 12; i++ {
		for j := 0; j < 3; j++ {
			simplex.grad3[i][j] = rand.Intn(3) - 1 // random values between -1 and 1
		}
	}

	p := rand.Perm(256) // generates a permutation of numbers from 0 to 255
	for i, v := range p {
		simplex.p[i] = v
		simplex.perm[i] = v
		simplex.perm[i+256] = v
	}

	return simplex
}
