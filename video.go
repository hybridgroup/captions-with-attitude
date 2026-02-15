package main

import (
	"errors"
	"image"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/mtmd"
)

// imgToBitmap converts an image.Image to an mtmd.Bitmap.
// It locks a mutex to ensure thread safety when accessing the image data.
func imgToBitmap(img image.Image) (mtmd.Bitmap, error) {
	if img == nil {
		return mtmd.Bitmap(0), errors.New("empty image")
	}

	mutex.Lock()
	defer mutex.Unlock()

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	rgb := make([]uint8, 0, width*height*3)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert from 16-bit to 8-bit color
			rgb = append(rgb, uint8(r>>8), uint8(g>>8), uint8(b>>8))
		}
	}

	bitmap := mtmd.BitmapInit(
		uint32(width),
		uint32(height),
		uintptr(unsafe.Pointer(&rgb[0])),
	)
	return bitmap, nil
}
