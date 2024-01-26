package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	perlin "github.com/CGuilhem/Perlin-Noise/Go/perlin/perlin"
)

func main() {
	typePtr := flag.String("type", "noise", "Parameter type")
	flag.Parse()

	start := time.Now()
	rand.Seed(time.Now().UnixNano())

	gradients := make([][][]float64, perlin.GridSize)
	for i := range gradients {
		gradients[i] = make([][]float64, perlin.GridSize)
		for j := range gradients[i] {
			gradients[i][j] = make([]float64, 2)
			angle := rand.Float64() * 2.0 * math.Pi
			gradients[i][j][0] = math.Cos(angle)
			gradients[i][j][1] = math.Sin(angle)
		}
	}

	if *typePtr == "map" {
		img := image.NewRGBA(image.Rect(0, 0, perlin.ImageWidth, perlin.ImageHeight))

		var wg sync.WaitGroup
		for i := 0; i < perlin.NumGoroutines; i++ {
			wg.Add(1)
			go perlin.GenerateImageRow(img, i, gradients, *typePtr, &wg)
		}

		wg.Wait()

		file, _ := os.Create("perlin_noise_MAP.jpeg")
		defer file.Close()
		jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	} else {
		img := image.NewGray(image.Rect(0, 0, perlin.ImageWidth, perlin.ImageHeight))

		var wg sync.WaitGroup
		for i := 0; i < perlin.NumGoroutines; i++ {
			wg.Add(1)
			go perlin.GenerateGrayImageRow(img, i, gradients, &wg)
		}

		wg.Wait()

		file, _ := os.Create("perlin_noise.jpeg")
		defer file.Close()
		jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	}

	fmt.Println("Time elapsed: ", time.Since(start))
}
