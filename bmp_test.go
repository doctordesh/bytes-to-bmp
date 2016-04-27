package main

import "testing"

func TestReverse(t *testing.T) {
	bmp := BMP{}
	bmp.data = []byte{0x00, 0x11, 0x22, 0x33}

	data := bmp.dataInReverse()

	if data[0] != 0x33 || data[1] != 0x22 || data[2] != 0x11 || data[3] != 0x00 {
		t.Error("Should be reverse")
	}
}

func TestAddRGB(t *testing.T) {
	bmp := BMP{Width: 1, Height: 1}
	bmp.Add(0x00, 0x11, 0x22)

	if bmp.data[0] != 0x22 || bmp.data[1] != 0x11 || bmp.data[2] != 0x00 {
		t.Error("Adding RGB does not reverse byte order")
	}
}
