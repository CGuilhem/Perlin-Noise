package main

import (
	"fmt"
	"image/color"

	simplex "github.com/CGuilhem/Perlin-Noise/Go/simplex/simplex"
	"github.com/fogleman/gg"
)

func main() {
	width := 512
	height := 512

	// Create a new context
	dc := gg.NewContext(width, height)

	// Create a Simplex instance
	simplex := simplex.NewSimplex()

	// Draw pixels based on Simplex noise
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Map x and y to a range suitable for Simplex noise
			scale := 0.02
			n := simplex.Noise(float64(x)*scale, float64(y)*scale)
			fmt.Println(n)

			// Map the noise value to a grayscale color
			colorValue := uint8((n + 1.0) * 0.5 * 255)
			pixelColor := color.Gray{Y: colorValue}

			// Set the pixel color
			dc.SetRGB(float64(pixelColor.Y)/255, float64(pixelColor.Y)/255, float64(pixelColor.Y)/255)

			// Draw a pixel at (x, y)
			dc.DrawPoint(float64(x), float64(y), 1)
			dc.Stroke()
		}
	}

	// Save the image to a file
	if err := dc.SavePNG("simplex_noise.png"); err != nil {
		fmt.Println("Error saving PNG:", err)
	} else {
		fmt.Println("Image saved successfully.")
	}
}
