package perlin

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// Fonction pour une transition de 0.0 à 1.0 dans la plage [0.0, 1.0]
func smoothstep(w float64) float64 {
	if w <= 0.0 {
		return 0.0
	}
	if w >= 1.0 {
		return 1.0
	}
	return w * w * (3.0 - 2.0*w)
}

// Fonction pour interpoler en douceur entre a0 et a1
// Le poids w devrait être dans la plage [0.0, 1.0]
func interpolate(a0, a1, w float64) float64 {
	return a0 + (a1-a0)*smoothstep(w)
}

// Calcule le produit scalaire des vecteurs de distance et de gradient.
func dotGridGradient(ix, iy int, x, y float64, gradients [][][]float64) float64 {
	// Vérifiez que les indices sont dans les limites du tableau gradients
	if iy >= 0 && iy < len(gradients) && ix >= 0 && ix < len(gradients[0]) {
		// Calcule le vecteur de distance
		dx := x - float64(ix)
		dy := y - float64(iy)

		// Vérifiez que les indices sont dans les limites du tableau de gradients[iy][ix]
		if len(gradients[iy]) > 0 && len(gradients[iy][ix]) >= 2 {
			// Calcule le produit scalaire
			return dx*gradients[iy][ix][0] + dy*gradients[iy][ix][1]
		}
	}

	// Si les indices sont hors limites, retournez une valeur par défaut
	return 0.0
}

// Calcule le bruit de Perlin aux coordonnées x, y
func Perlin(x, y float64, gradients [][][]float64) float64 {
	// Détermine les coordonnées de la cellule de la grille
	x0 := int(math.Floor(x))
	x1 := x0 + 1
	y0 := int(math.Floor(y))
	y1 := y0 + 1

	// Détermine les poids d'interpolation
	// On pourrait aussi utiliser une courbe polynomiale d'ordre supérieur ici
	sx := x - float64(x0)
	sy := y - float64(y0)

	// Interpolation entre les gradients des points de la grille
	var n0, n1, ix0, ix1, value float64
	n0 = dotGridGradient(x0, y0, x, y, gradients)
	n1 = dotGridGradient(x1, y0, x, y, gradients)
	ix0 = interpolate(n0, n1, sx)
	n0 = dotGridGradient(x0, y1, x, y, gradients)
	n1 = dotGridGradient(x1, y1, x, y, gradients)
	ix1 = interpolate(n0, n1, sx)
	value = interpolate(ix0, ix1, sy)

	return value
}

const (
	GridSize      = 256
	ImageWidth    = 800
	ImageHeight   = 800
	NumGoroutines = ImageHeight
)

func GenerateImageRow(img *image.RGBA, row int, gradients [][][]float64, typePtr string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < ImageWidth; i++ {
		x, y := float64(i)/float64(ImageWidth), float64(row)/float64(ImageHeight)
		frequency := 4.0
		x, y = x*frequency, y*frequency

		octaves := 10
		persistence := 0.85
		total := 0.0
		amplitude := 1.0
		for o := 0; o < octaves; o++ {
			total += Perlin(x, y, gradients) * amplitude
			x, y = x*2, y*2
			amplitude *= persistence
		}

		value := (total + 1.0) / 2.0

		var col color.Color
		if typePtr == "map" {
			col = GetColor(value)
		} else {
			gray := uint8((total + 1.0) / 2.0 * 255.0)
			col = color.Gray{Y: gray}
		}

		img.Set(i, row, col)
	}
}

func GenerateGrayImageRow(img *image.Gray, row int, gradients [][][]float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < ImageWidth; i++ {
		x, y := float64(i)/float64(ImageWidth), float64(row)/float64(ImageHeight)
		frequency := 4.0
		x, y = x*frequency, y*frequency

		octaves := 10
		persistence := 0.65
		total := 0.0
		amplitude := 1.0
		for o := 0; o < octaves; o++ {
			total += Perlin(x, y, gradients) * amplitude
			x, y = x*2, y*2
			amplitude *= persistence
		}

		gray := uint8((total + 1.0) / 2.0 * 255.0)
		img.SetGray(i, row, color.Gray{Y: gray})
	}
}

func GetColor(value float64) color.Color {
	switch {
	case value < 0.2:
		return color.RGBA{0, 0, 255, 255} // water
	case value < 0.3:
		return color.RGBA{194, 178, 128, 255} // sand
	case value < 0.5:
		return color.RGBA{34, 139, 34, 255} // forest
	case value < 0.7:
		return color.RGBA{139, 69, 19, 255} // mountain
	default:
		return color.RGBA{255, 250, 250, 255} // snow
	}
}
