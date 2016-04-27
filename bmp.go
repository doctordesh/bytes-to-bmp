package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func print(data []byte) {
	fmt.Printf("%#x\n", data)
}

func reverse(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	return append(reverse(data[1:]), data[0])
}

type BMP struct {
	Width  uint
	Height uint
	data   []byte
}

func (bmp *BMP) append(v ...byte) {
	bmp.data = append(bmp.data, v...)
}

func (bmp *BMP) AddFromReader(r io.Reader) {
	rgb := make([]byte, 3)

	fmt.Printf("Width: %d, Height: %d, Total: %d\n", bmp.Width, bmp.Height, bmp.Width*bmp.Height)

	for i := 0; i < int(bmp.Width*bmp.Height); i++ {
		_, err := r.Read(rgb)
		check(err)
		bmp.Add(rgb[0], rgb[1], rgb[2])
	}
}

func (bmp *BMP) Add(r, g, b byte) {
	//bmp.append(r, g, b)
	bmp.append(b, g, r)
}

func (bmp *BMP) Out(filename string) {
	var d []byte
	var image []byte

	image = bmp.buildImage(bmp.data)

	d = append(d, 0x42, 0x4D)                           // 0  - 2  signature, must be 4D42 hex, ASCII 'B' and 'M'
	d = append(d, bmp.uint32tob(bmp.length())...)       // 2  - 4  size of BMP file in bytes (unreliable)
	d = append(d, 0x00, 0x00)                           // 6  - 2  reserved, must be zero
	d = append(d, 0x00, 0x00)                           // 8  - 2  reserved, must be zero
	d = append(d, bmp.uint32tob(bmp.HeaderSize())...)   // 10 - 2  offset to start of image data in bytes
	d = append(d, bmp.uint32tob(40)...)                 // 14 - 4  size of BITMAPINFOHEADER structure, must be 40
	d = append(d, bmp.uint32tob(uint32(bmp.Width))...)  // 18 - 4  image width in pixels
	d = append(d, bmp.uint32tob(uint32(bmp.Height))...) // 22 - 4  image height in pixels
	d = append(d, bmp.uint16tob(1)...)                  // 26 - 4  number of planes in the image, must be 1
	d = append(d, bmp.uint16tob(24)...)                 // 28 - 2  number of bits per pixel (1, 4, 8, or 24)
	d = append(d, bmp.uint32tob(0)...)                  // 30 - 2  compression type (0=none, 1=RLE-8, 2=RLE-4)
	d = append(d, bmp.uint32tob(uint32(len(image)))...) // 34 - 4  size of image data in bytes (including padding)
	d = append(d, bmp.uint32tob(4875)...)               // 38 - 4  horizontal resolution in pixels per meter (unreliable)
	d = append(d, bmp.uint32tob(4875)...)               // 42 - 4  vertical resolution in pixels per meter (unreliable)
	d = append(d, bmp.uint32tob(0)...)                  // 46 - 4  number of colors in image, or zero
	d = append(d, bmp.uint32tob(0)...)                  // 50 - 4  number of important colors, or zero

	d = append(d, image...)

	ioutil.WriteFile(filename, d, os.ModePerm)
}

func (bmp *BMP) buildImage(data []byte) []byte {
	var rowLength, padding uint
	var image []byte

	rowLength = bmp.Width * 3
	padding = 4 - rowLength%4

	if padding == 4 {
		padding = 0
	}

	for y := uint(0); y < bmp.Height; y++ {
		image = append(image, data[len(data)-int(rowLength):]...)
		data = data[:len(data)-int(rowLength)]
		for i := 0; uint(i) < padding; i++ {
			image = append(image, 0x00)
		}
	}

	return image

	// var x, y, i, rowLength uint
	// var image []byte

	// rowLength = bmp.Width * 3 // 3 bytes per pixel

	// for y = 0; y < bmp.Height; y++ {
	// 	for x = 0; x < rowLength; x++ {
	// 		image = append(image, data[i])
	// 		i++
	// 	}

	// 	for i = 0; i < (rowLength % 4); i++ {
	// 		fmt.Printf("padding\n")
	// 		image = append(image, 0x00)
	// 	}
	// }

	// fmt.Printf("Image: %#x\n", image)

	// return image
}

func (bmp *BMP) uint32tob(i uint32) []byte {
	var b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

func (bmp *BMP) uint16tob(i uint16) []byte {
	var b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

func (bmp *BMP) dataLength() uint32 {
	return uint32(len(bmp.data))
}

func (bmp *BMP) length() uint32 {
	return bmp.dataLength() + bmp.HeaderSize()
}

func (bmp *BMP) dataInReverse() []byte {
	var d []byte
	for i := len(bmp.data); i > 0; i-- {
		d = append(d, bmp.data[i-1])
	}
	return d
}

func (bmp *BMP) HeaderSize() uint32 {
	return 54
}

var data []byte

func add(v ...byte) {
	data = append(data, v...)
}

func addRGB(r, g, b byte) {
	data = append(data, b, g, r)
}

//func main() {
// add(0x42, 0x4D)             // 0      signature, must be 4D42 hex
// add(0x3E, 0x00, 0x00, 0x00) // 2      size of BMP file in bytes (unreliable)
// add(0x00, 0x00)             // 6      reserved, must be zero
// add(0x00, 0x00)             // 8      reserved, must be zero
// add(0x36, 0x00, 0x00, 0x00) // 10     offset to start of image data in bytes
// add(0x28, 0x00, 0x00, 0x00) // 14     size of BITMAPINFOHEADER structure, must be 40
// add(0x01, 0x00, 0x00, 0x00) // 18     image width in pixels
// add(0x02, 0x00, 0x00, 0x00) // 22     image height in pixels
// add(0x01, 0x00)             // 26     number of planes in the image, must be 1
// add(0x18, 0x00)             // 28     number of bits per pixel (1, 4, 8, or 24)
// add(0x00, 0x00, 0x00, 0x00) // 30     compression type (0=none, 1=RLE-8, 2=RLE-4)
// add(0x10, 0x00, 0x00, 0x00) // 34     size of image data in bytes (including padding)
// add(0x13, 0x0B, 0x00, 0x00) // 38     horizontal resolution in pixels per meter (unreliable)
// add(0x13, 0x0B, 0x00, 0x00) // 42     vertical resolution in pixels per meter (unreliable)
// add(0x00, 0x00, 0x00, 0x00) // 46     number of colors in image, or zero
// add(0x00, 0x00, 0x00, 0x00) // 50     number of important colors, or zero

// addRGB(0xFF, 0x00, 0x00)
// add(0x00) // padding
// addRGB(0x00, 0xFF, 0x00)
// add(0x00) // padding
// add(0x00, 0xFF, 0x00)
// add(0x00, 0x00, 0xFF)
// add(0x00, 0xFF, 0xFF)

// add(0xFF, 0x00, 0x00)
// add(0x00, 0xFF, 0xFF)
// add(0x00) // padding
// add(0x00) // padding

// ioutil.WriteFile("tmp2.bmp", data, os.ModePerm)

//}
