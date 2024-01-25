package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

// Fonction pour une transition en douceur de 0.0 à 1.0 dans la plage [0.0, 1.0]
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
func perlin(x, y float64, gradients [][][]float64) float64 {
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

// Définissez la taille de votre grille de gradients
const gridSize = 256

func main() {

	typePtr := flag.String("type", "noise", "Parameter type")
	flag.Parse()

	start := time.Now()

	// Initialisez le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Créez un tableau 3D pour stocker les gradients
	gradients := make([][][]float64, gridSize)
	for i := range gradients {
		gradients[i] = make([][]float64, gridSize)
		for j := range gradients[i] {
			gradients[i][j] = make([]float64, 2) // Nous utilisons des gradients 2D pour le bruit de Perlin 2D

			// Générez un vecteur aléatoire
			angle := rand.Float64() * 2.0 * math.Pi // Angle aléatoire
			gradients[i][j][0] = math.Cos(angle)    // Composante x du gradient
			gradients[i][j][1] = math.Sin(angle)    // Composante y du gradient
		}
	}

	if *typePtr == "map" {
		// Créez une nouvelle image
		img := image.NewRGBA(image.Rect(0, 0, 800, 800))

		// Définissez les couleurs pour chaque type de terrain
		water := color.RGBA{0, 0, 255, 255}      // Bleu pour l'eau
		sand := color.RGBA{194, 178, 128, 255}   // Beige pour le sable
		forest := color.RGBA{34, 139, 34, 255}   // Vert pour la forêt
		mountain := color.RGBA{139, 69, 19, 255} // Marron pour la montagne
		snow := color.RGBA{255, 250, 250, 255}   // Blanc pour la neige

		// Remplissez l'image avec des valeurs de bruit de Perlin
		for i := 0; i < img.Bounds().Dx(); i++ {
			for j := 0; j < img.Bounds().Dy(); j++ {
				// Générez une valeur de bruit de Perlin pour ce point
				x, y := float64(i)/800.0, float64(j)/800.0 // Coordonnées normalisées
				frequency := 5.0                           // Augmentez la fréquence pour un bruit plus "fin"
				x, y = x*frequency, y*frequency            // Ajustez les coordonnées en fonction de la fréquence

				// Ajoutez plusieurs octaves de bruit
				octaves := 10
				persistence := 0.85
				total := 0.0
				amplitude := 1.0
				for o := 0; o < octaves; o++ {
					total += perlin(x, y, gradients) * amplitude

					// Préparez la prochaine octave
					x, y = x*2, y*2          // Doublez la fréquence
					amplitude *= persistence // Réduisez l'amplitude
				}
				// Convertissez le bruit (qui est entre -1 et 1) en une valeur entre 0 et 1
				value := (total + 1.0) / 2.0

				// Définissez la couleur de ce pixel dans l'image en fonction de la valeur du bruit
				var col color.RGBA
				switch {
				case value < 0.2:
					col = water
				case value < 0.3:
					col = sand
				case value < 0.5:
					col = forest
				case value < 0.7:
					col = mountain
				default:
					col = snow
				}
				img.Set(i, j, col)
			}
		}

		// Écrivez l'image dans un fichier
		file, _ := os.Create("perlin_noise_MAP.png")
		defer file.Close()
		png.Encode(file, img)
	} else {
		// Créez une nouvelle image
		img := image.NewGray(image.Rect(0, 0, 800, 800))

		// Remplissez l'image avec des valeurs de bruit de Perlin
		for i := 0; i < img.Bounds().Dx(); i++ {
			for j := 0; j < img.Bounds().Dy(); j++ {
				// Générez une valeur de bruit de Perlin pour ce point
				x, y := float64(i)/800.0, float64(j)/800.0 // Coordonnées normalisées
				frequency := 4.0                           // Augmentez la fréquence pour un bruit plus "fin"
				x, y = x*frequency, y*frequency            // Ajustez les coordonnées en fonction de la fréquence

				// Ajoutez plusieurs octaves de bruit
				octaves := 6
				persistence := 0.5
				total := 0.0
				amplitude := 1.0
				for o := 0; o < octaves; o++ {
					total += perlin(x, y, gradients) * amplitude

					// Préparez la prochaine octave
					x, y = x*2, y*2          // Doublez la fréquence
					amplitude *= persistence // Réduisez l'amplitude
				}

				// Convertissez le bruit (qui est entre -1 et 1) en une valeur de gris (entre 0 et 255)
				gray := uint8((total + 1.0) / 2.0 * 255.0)

				// Définissez la couleur de ce pixel dans l'image
				img.SetGray(i, j, color.Gray{Y: gray})
			}
		}

		// Écrivez l'image dans un fichier
		file, _ := os.Create("perlin_noise.png")
		defer file.Close()
		png.Encode(file, img)
	}

	fmt.Println("Time elapsed: ", time.Since(start))

}
