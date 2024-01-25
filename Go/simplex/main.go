package main

import (
	"flag"
	"fmt"
	"image/color"
	"time"

	simplex "github.com/CGuilhem/Perlin-Noise/Go/simplex/simplex"
	"github.com/fogleman/gg"
)

func main() {

	typePtr := flag.String("type", "noise", "Parameter type")

	flag.Parse()
	start := time.Now()

	width := 512
	height := 512

	dc := gg.NewContext(width, height)

	simplex := simplex.NewSimplex()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			scale := 0.02
			n := simplex.Noise(float64(x)*scale, float64(y)*scale)

			pixelColor := getColor(n, typePtr)
			dc.SetColor(pixelColor)
			dc.DrawPoint(float64(x), float64(y), 1)
			dc.Stroke()
		}
	}

	if err := dc.SavePNG("simplex_noise.png"); err != nil {
		fmt.Println("Error saving PNG:", err)
	} else {
		fmt.Println("Image saved successfully.")
	}

	fmt.Println("Time elapsed:", time.Since(start))
}

func getColor(n float64, typePtr *string) color.Color {

	if *typePtr == "map" {
		elevation := (n + 1.0) * 0.5

		switch {
		case elevation < 0.4:
			return color.RGBA{0, 0, 255, 255} // water
		case elevation < 0.5:
			return color.RGBA{210, 180, 140, 255} // sand
		case elevation < 0.6:
			return color.RGBA{34, 139, 34, 255} // forest
		case elevation < 0.8:
			return color.RGBA{128, 128, 128, 255} // rock mountain
		default:
			return color.RGBA{255, 255, 255, 255} // snow
		}
	} else {
		colorValue := uint8((n + 1.0) * 0.5 * 255)
		return color.Gray{Y: colorValue}
	}
}
