package main

import (
	"math"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getSide(dataLength int) int {
	return int(math.Sqrt(float64(dataLength / 3)))
}

func main() {
	var err error
	var bmp BMP
	// bmp1 := BMP{}
	// bmp1.Width = 3
	// bmp1.Height = 3
	// bmp1.Add(0x00, 0x00, 0x00) // Black
	// bmp1.Add(0x77, 0x77, 0x77) // Gray
	// bmp1.Add(0xFF, 0xFF, 0xFF) // White
	// bmp1.Add(0xFF, 0xFF, 0x00) // Yellow
	// bmp1.Add(0x00, 0xFF, 0xFF) // Cyan
	// bmp1.Add(0xFF, 0x00, 0xFF) // Magenta
	// bmp1.Add(0xFF, 0x00, 0x00) // Red
	// bmp1.Add(0x00, 0xFF, 0x00) // Green
	// bmp1.Add(0x00, 0x00, 0xFF) // Blue

	// bmp1.Out("bmp1.bmp")

	args := os.Args[1:]

	f, err := os.Open(args[0])
	check(err)

	info, err := f.Stat()

	bmp = BMP{}
	bmp.Width = uint(getSide(int(info.Size())))
	bmp.Height = bmp.Width
	bmp.AddFromReader(f)
	bmp.Out("image.bmp")
}
