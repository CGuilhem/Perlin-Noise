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
	"sync"
	"time"

	"github.com/CGuilhem/Perlin-Noise/Go/perlin/perlin"
)

const (
	gridSize      = 256
	imageWidth    = 800
	imageHeight   = 800
	numGoroutines = imageHeight
)

func generateImageRow(img *image.RGBA, row int, gradients [][][]float64, typePtr string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < imageWidth; i++ {
		x, y := float64(i)/float64(imageWidth), float64(row)/float64(imageHeight)
		frequency := 4.0
		x, y = x*frequency, y*frequency

		octaves := 10
		persistence := 0.85
		total := 0.0
		amplitude := 1.0
		for o := 0; o < octaves; o++ {
			total += perlin.Perlin(x, y, gradients) * amplitude
			x, y = x*2, y*2
			amplitude *= persistence
		}

		value := (total + 1.0) / 2.0

		var col color.Color
		if typePtr == "map" {
			col = getColor(value)
		} else {
			gray := uint8((total + 1.0) / 2.0 * 255.0)
			col = color.Gray{Y: gray}
		}

		img.Set(i, row, col)
	}
}

func generateGrayImageRow(img *image.Gray, row int, gradients [][][]float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < imageWidth; i++ {
		x, y := float64(i)/float64(imageWidth), float64(row)/float64(imageHeight)
		frequency := 4.0
		x, y = x*frequency, y*frequency

		octaves := 10
		persistence := 0.65
		total := 0.0
		amplitude := 1.0
		for o := 0; o < octaves; o++ {
			total += perlin.Perlin(x, y, gradients) * amplitude
			x, y = x*2, y*2
			amplitude *= persistence
		}

		gray := uint8((total + 1.0) / 2.0 * 255.0)
		img.SetGray(i, row, color.Gray{Y: gray})
	}
}

func getColor(value float64) color.Color {
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

func main() {
	typePtr := flag.String("type", "noise", "Parameter type")
	flag.Parse()

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

	if *typePtr == "map" {
		img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go generateImageRow(img, i, gradients, *typePtr, &wg)
		}

		wg.Wait()

		file, _ := os.Create("perlin_noise_MAP.png")
		defer file.Close()
		png.Encode(file, img)
	} else {
		img := image.NewGray(image.Rect(0, 0, imageWidth, imageHeight))

		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go generateGrayImageRow(img, i, gradients, &wg)
		}

		wg.Wait()

		file, _ := os.Create("perlin_noise.png")
		defer file.Close()
		png.Encode(file, img)
	}

	fmt.Println("Time elapsed: ", time.Since(start))
}
