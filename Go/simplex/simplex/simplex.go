package simplex

type Simplex struct {
	grad3 [12][3]int
	p     [256]int
	perm  [512]int
}

func (s *Simplex) dot(g [3]int, x, y float64) float64 {
	return float64(g[0])*x + float64(g[1])*y
}

func (s *Simplex) noise(xin, yin float64) {
	// F2 := 0.5 * (math.Sqrt(3.0) - 1.0)
	// s := (xin + yin) * F2
	// i := math.Floor(xin + s)
	// j := math.Floor(yin + s)

	// G2 := (3.0 - math.Sqrt(3.0)) / 6.0
	// t := (i + j) * G2
	// X0 := i - t
	// Y0 := j - t
	// x0 := xin - X0
	// y0 := yin - Y0

	// var i1, j1 int
	// if x0 > y0 {
	// 	i1, j1 = 1, 0
	// } else {
	// 	i1, j1 = 0, 1
	// }

	// x1 := x0 - float64(i1) + G2
	// y1 := y0 - float64(j1) + G2
	// x2 := x0 - 1.0 + 2.0*G2
	// y2 := y0 - 1.0 + 2.0*G2

	// ii := int(i) & 255
	// jj := int(j) & 255
	// gi0 := s.perm[ii+s.perm[jj]] % 12
	// gi1 := s.perm[ii+i1+s.perm[jj+j1]] % 12
	// gi2 := s.perm[ii+1+s.perm[jj+1]] % 12

	// t0 := 0.5 - x0*x0 - y0*y0
	// var n0 float64
	// if t0 < 0 {
	// 	n0 = 0.0
	// } else {
	// 	t0 *= t0
	// 	n0 = t0 * t0 * s.dot(s.grad3[gi0], x0, y0)
	// }

	// t1 := 0.5 - x1*x1 - y1*y1
	// var n1 float64
	// if t1 < 0 {
	// 	n1 = 0.0
	// } else {
	// 	t1 *= t1
	// 	n1 = t1 * t1 * s.dot(s.grad3[gi1], x1, y1)
	// }

	// t2 := 0.5 - x2*x2 - y2*y2
	// var n2 float64
	// if t2 < 0 {
	// 	n2 = 0.0
	// } else {
	// 	t2 *= t2
	// 	n2 = t2 * t2 * s.dot(s.grad3[gi2], x2, y2)
	// }

	// return 70.0 * (n0 + n1 + n2)
}

func NewSimplex() *Simplex {
	simplex := &Simplex{}
	for x := 0; x < 255; x++ {
		simplex.perm[x] = simplex.p[x]
		simplex.perm[x+256] = simplex.p[x]
	}
	return simplex
}
