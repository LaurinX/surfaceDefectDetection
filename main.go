package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	// Set parameters for the image
	width, height := 10, 10
	pointX, pointY := 9, 0
	mmToPixelRatio := 0.5 // example: 0.5 mm per pixel
	diameterMM := 8.0     // cylinder diameter in mm
	rotationAngle := 45.0 // rotation angle in degrees

	// Create the image
	img := CreateImageWithPoint(width, height, pointX, pointY)

	// Save the image as a PNG file
	err := SaveImage(img, "output.png")
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Println("Image saved successfully as output.png")

	// Rotate and Shift 10 times
	for i := 0; i < 10; i++ {
		// Rotate the image
		rotatedImg := RotateCylinderImage(img, mmToPixelRatio, diameterMM, rotationAngle)

		// Save the rotated image
		err = SaveImage(rotatedImg, "rotated_output_"+strconv.Itoa(i)+".png")
		if err != nil {
			log.Fatalf("Failed to save rotated image: %v", err)
		}
		img = rotatedImg
	}

	log.Println("Rotated image saved successfully as rotated_output.png")
}

// CreateImageWithPoint creates a grayscale image with a specified size and a black point at (pointX, pointY)
func CreateImageWithPoint(width, height, pointX, pointY int) *image.Gray {
	if !(pointX < width && pointY < height && pointX >= 0 && pointY >= 0) {
		log.Fatal("Warning: Point", pointX, pointY, " is outside the image bounds", pointX, pointY)
	}
	// Create a new grayscale image
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Set the background to white
	white := color.Gray{Y: 230}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, white)
		}
	}

	// Set the specified point to black
	black := color.Gray{Y: 0}

	img.Set(pointX, pointY, black)

	return img
}

// SaveImage saves an image to a specified file path
func SaveImage(img *image.Gray, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// RotateCylinderImage takes an image, mm-to-pixel ratio, diameter in mm, and rotation angle in degrees, then transforms the image accordingly.
func RotateCylinderImage(img *image.Gray, mmToPixelRatio float64, diameterMM float64, rotationAngle float64) *image.Gray {
	// Calculate the circumference in mm
	circumferenceMM := math.Pi * diameterMM

	// Convert rotation angle to a pixel shift
	rotationRadians := rotationAngle * math.Pi / 180.0                // Convert to radians
	rotationMM := rotationRadians * (circumferenceMM / (2 * math.Pi)) // Distance traveled in mm
	pixelShift := int(rotationMM * mmToPixelRatio)
	log.Println(pixelShift, rotationMM*mmToPixelRatio) // Distance in pixels

	// Prepare a new blank image with the same size as the input
	rotatedImg := image.NewGray(img.Bounds())

	// Shift pixels in the original image horizontally based on the calculated pixelShift
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Wrap around each row to simulate cylindrical rotation
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the new x position with wrapping
			newX := (x + pixelShift) % width
			if newX < 0 {
				newX += width
			}
			rotatedImg.SetGray(newX, y, img.GrayAt(x, y))
		}
	}

	return rotatedImg
}
