package parse

import (
	"image"
	"image/color"
	_ "image/png" // register PNG decoder
	"math"
	"os"
)

// dBZColor represents a color and its corresponding dBZ value.
type dBZColor struct {
	Color color.Color
	dBZ   int
}

// dBZScale is the color scale for the IMD Doppler radar.
var dBZScale = []dBZColor{
	{Color: color.RGBA{0, 0, 255, 255}, dBZ: 25},    // Light Blue
	{Color: color.RGBA{0, 255, 0, 255}, dBZ: 35},    // Green
	{Color: color.RGBA{255, 255, 0, 255}, dBZ: 45},   // Yellow
	{Color: color.RGBA{255, 0, 0, 255}, dBZ: 55},    // Red
	{Color: color.RGBA{128, 0, 128, 255}, dBZ: 65},  // Purple
	{Color: color.RGBA{255, 255, 255, 255}, dBZ: 70}, // White
}

// findClosestDBZ finds the closest dBZ value for a given color.
func findClosestDBZ(c color.Color) int {
	minDist := math.MaxFloat64
	maxDBZ := 0

	for _, scaleColor := range dBZScale {
		r1, g1, b1, _ := c.RGBA()
		r2, g2, b2, _ := scaleColor.Color.RGBA()

		dist := math.Sqrt(math.Pow(float64(r1)-float64(r2), 2) + math.Pow(float64(g1)-float64(g2), 2) + math.Pow(float64(b1)-float64(b2), 2))

		if dist < minDist {
			minDist = dist
			maxDBZ = scaleColor.dBZ
		}
	}
	return maxDBZ
}

// ParseRadarImage analyzes a radar image and returns the maximum dBZ value
// within a given radius of the center.
func ParseRadarImage(imagePath string, radiusKm int) (int, error) {
	f, err := os.Open(imagePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return 0, err
	}

	bounds := img.Bounds()
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	// Assuming a 250km range for the radar image.
	// This is a common range for IMD radars.
	// We can make this configurable later if needed.
	kmPerPixel := 250.0 / float64(bounds.Dx()/2)
	radiusPixels := int(float64(radiusKm) / kmPerPixel)

	maxDBZ := 0

	for y := centerY - radiusPixels; y <= centerY+radiusPixels; y++ {
		for x := centerX - radiusPixels; x <= centerX+radiusPixels; x++ {
			dist := math.Sqrt(math.Pow(float64(x-centerX), 2) + math.Pow(float64(y-centerY), 2))
			if dist <= float64(radiusPixels) {
				c := img.At(x, y)
				dbz := findClosestDBZ(c)
				if dbz > maxDBZ {
					maxDBZ = dbz
				}
			}
		}
	}

	return maxDBZ, nil
}
