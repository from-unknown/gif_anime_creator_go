package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"

	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/soniakeys/quant/median"
)

func main() {
	// command line arguments
	var argStartFlag string
	var argEndFlag string
	var argLoopFlag string
	// starting position
	flag.StringVar(&argStartFlag, "start", "left", "Please set direction for gif animation.")
	// ending position
	flag.StringVar(&argEndFlag, "end", "right", "Please set direction for gif animation.")
	// loop flag
	flag.StringVar(&argLoopFlag, "loop", "false", "Please set direction for gif animation.")
	flag.Parse()

	const postFix string = "_animated"

	// image file path
	filePath := flag.Args()[0]
	base := filepath.Base(filePath)
	ext := filepath.Ext(filePath)
	ext = strings.ToLower(ext)

	imageFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	var decodedImage image.Image

	// decide file type from file extention
	if ext == ".jpg" || ext == ".jpeg" {
		decodedImage, err = jpeg.Decode(imageFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if ext == ".png" {
		decodedImage, err = png.Decode(imageFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if ext == ".gif" {
		decodedImage, err = gif.Decode(imageFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		return
	}

	// open file to encode gif file
	tmpFileName := base[0:len(base)-len(ext)] + postFix + ".gif"
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer tmpFile.Close()

	var gifImage []*image.Paletted
	var delays []int
	var disposal []byte
	rect := decodedImage.Bounds()
	for i := 0; i < 10; i++ {
		var drawXBound int
		var drawYBound int
		drawXBound = 0
		drawYBound = 0
		if i < 5 { // for first 5 images
			if argLoopFlag == "true" { // if true, start from center
				switch argStartFlag {
				case "right":
					drawXBound = -rect.Dx() + int(math.Floor(float64(rect.Dx())*0.2*float64(i+5)))
				case "left":
					drawXBound = rect.Dx() - int(math.Floor(float64(rect.Dx())*0.2*float64(i+5)))
				case "top":
					drawYBound = rect.Dy() - int(math.Floor(float64(rect.Dy())*0.2*float64(i+5)))
				case "bottom":
					drawYBound = -rect.Dy() + int(math.Floor(float64(rect.Dy())*0.2*float64(i+5)))
				default:
				}
			} else {
				switch argStartFlag {
				case "right":
					drawXBound = -rect.Dx() + int(math.Floor(float64(rect.Dx())*0.2*float64(i)))
				case "left":
					drawXBound = rect.Dx() - int(math.Floor(float64(rect.Dx())*0.2*float64(i)))
				case "top":
					drawYBound = rect.Dy() - int(math.Floor(float64(rect.Dy())*0.2*float64(i)))
				case "bottom":
					drawYBound = -rect.Dy() + int(math.Floor(float64(rect.Dy())*0.2*float64(i)))
				default:
				}
			}
		} else { // for last 5 images
			if argLoopFlag == "true" {
				switch argEndFlag {
				case "right":
					drawXBound = -int(math.Floor(float64(rect.Dx()) * 0.2 * float64(i-10)))
				case "left":
					drawXBound = int(math.Floor(float64(rect.Dx()) * 0.2 * float64(i-10)))
				case "top":
					drawYBound = int(math.Floor(float64(rect.Dy()) * 0.2 * float64(i-10)))
				case "bottom":
					drawYBound = -int(math.Floor(float64(rect.Dy()) * 0.2 * float64(i-10)))
				default:
				}
			} else {
				switch argEndFlag {
				case "right":
					drawXBound = -int(math.Floor(float64(rect.Dx()) * 0.2 * float64(i-5)))
				case "left":
					drawXBound = int(math.Floor(float64(rect.Dx()) * 0.2 * float64(i-5)))
				case "top":
					drawYBound = int(math.Floor(float64(rect.Dy()) * 0.2 * float64(i-5)))
				case "bottom":
					drawYBound = -int(math.Floor(float64(rect.Dy()) * 0.2 * float64(i-5)))
				default:
				}
			}
		}
		drawPoint := image.Pt(drawXBound, drawYBound)

		// gif can have 256 colors and one for Transparent
		const n = 255

		// use Quantizer to choose colors for gif
		var q draw.Quantizer = median.Quantizer(n)
		quantizePallete := q.Quantize(make(color.Palette, 0, n), decodedImage)

		newPalette := color.Palette{
			image.Transparent,
		}
		for _, color := range quantizePallete {
			newPalette = append(newPalette, color)
		}
		tmpPalette := image.NewPaletted(rect, newPalette)
		draw.Draw(tmpPalette, rect, decodedImage, drawPoint, draw.Src)
		// if use FloydSteinberg, some images output become wired color
		//draw.FloydSteinberg.Draw(tmpPalette, rect, decodedImage, drawPoint)

		gifImage = append(gifImage, tmpPalette)
		delays = append(delays, 0)
		disposal = append(disposal, 2)
	}

	// Set size to resized size
	gif.EncodeAll(tmpFile, &gif.GIF{
		Image:     gifImage,
		Delay:     delays,
		Disposal:  disposal,
		LoopCount: 0,
	})
}

// Check if color is already in the Palette
func contains(colorPalette color.Palette, c color.Color) bool {
	for _, tmpColor := range colorPalette {
		if tmpColor == c {
			return true
		}
	}
	return false
}
