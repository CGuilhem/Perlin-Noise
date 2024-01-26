package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	perlin "github.com/CGuilhem/Perlin-Noise/Go/perlin/perlin"
)

const (
	gridSize      = 256
	imageWidth    = 800
	imageHeight   = 800
	numGoroutines = imageHeight
)

func main() {

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		// Définir les en-têtes pour désactiver CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		fmt.Println("GET /image")

		start := time.Now()

		rand.Seed(time.Now().UnixNano())

		gradients := make([][][]float64, gridSize)
		for i := range gradients {
			gradients[i] = make([][]float64, gridSize)
			for j := range gradients[i] {
				gradients[i][j] = make([]float64, 2)
				angle := rand.Float64() * 2.0 * math.Pi
				gradients[i][j][0] = math.Cos(angle)
				gradients[i][j][1] = math.Sin(angle)
			}
		}

		img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go perlin.GenerateImageRow(img, i, gradients, "map", &wg)
		}

		wg.Wait()

		file, _ := os.Create("perlin_noise_MAP.jpeg")
		defer file.Close()
		jpeg.Encode(file, img, &jpeg.Options{Quality: 80})

		fmt.Println("Image generated in", time.Since(start))

		http.ServeFile(w, r, "perlin_noise_MAP.jpeg")
	})

	http.ListenAndServe(":8080", nil)
}
