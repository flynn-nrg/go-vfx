package oiio

/*
#cgo CXXFLAGS: -std=c++17
#cgo pkg-config: OpenImageIO fmt
#include <stdlib.h>
#include "./oiio_wrapper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"math"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/flynn-nrg/floatimage/floatimage"
)

func ReadImage32(filename string) (*floatimage.Float32NRGBA, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var errorMsg *C.char
	cImage := C.read_image(cFilename, &errorMsg)
	if cImage == nil {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return nil, fmt.Errorf("failed to read image: %s", err)
		}
		return nil, fmt.Errorf("failed to read image")
	}
	defer C.free_image(cImage)

	width := int(cImage.width)
	height := int(cImage.height)
	numChannels := int(cImage.channels)
	cData := (*[1 << 30]C.float)(unsafe.Pointer(cImage.data))[: width*height*numChannels : width*height*numChannels]

	var data []float32

	switch cImage.channels {
	case 3:
		data = toRGB32Slice(cData, width, height)
	case 4:
		data = toRGBA32Slice(cData, width, height, numChannels)
	default:
		return nil, fmt.Errorf("unsupported number of channels: %d", cImage.channels)
	}

	return floatimage.NewFloat32NRGBA(image.Rect(0, 0, width, height), data), nil
}

// ReadImage32Aces reads an image and converts it to ACEScg (AP1) color space using OpenColorIO.
// This function uses OIIO's built-in OCIO integration to:
// 1. Interpret the source image as "Utility - sRGB - Texture" (linearizes and applies correct primaries)
// 2. Convert to ACEScg working color space
// Note: Requires an OCIO configuration to be available (via OCIO env var or OIIO defaults)
func ReadImage32Aces(filename string) (*floatimage.Float32NRGBA, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var errorMsg *C.char
	cImage := C.read_image_aces(cFilename, &errorMsg)
	if cImage == nil {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return nil, fmt.Errorf("failed to read image for ACEScg: %s", err)
		}
		return nil, fmt.Errorf("failed to read image for ACEScg")
	}
	defer C.free_image(cImage)

	width := int(cImage.width)
	height := int(cImage.height)
	numChannels := int(cImage.channels)
	cData := (*[1 << 30]C.float)(unsafe.Pointer(cImage.data))[: width*height*numChannels : width*height*numChannels]

	var data []float32

	switch cImage.channels {
	case 3:
		data = toRGB32Slice(cData, width, height)
	case 4:
		data = toRGBA32Slice(cData, width, height, numChannels)
	default:
		return nil, fmt.Errorf("unsupported number of channels: %d", cImage.channels)
	}

	return floatimage.NewFloat32NRGBA(image.Rect(0, 0, width, height), data), nil
}

func ReadImage64(filename string) (*floatimage.Float64NRGBA, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var errorMsg *C.char
	cImage := C.read_image(cFilename, &errorMsg)
	if cImage == nil {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return nil, fmt.Errorf("failed to read image: %s", err)
		}
		return nil, fmt.Errorf("failed to read image")
	}
	defer C.free_image(cImage)

	width := int(cImage.width)
	height := int(cImage.height)
	numChannels := int(cImage.channels)
	cData := (*[1 << 30]C.float)(unsafe.Pointer(cImage.data))[: width*height*numChannels : width*height*numChannels]

	var data []float64

	switch cImage.channels {
	case 3:
		data = toRGB64Slice(cData, width, height)
	case 4:
		data = toRGBA64Slice(cData, width, height, numChannels)
	default:
		return nil, fmt.Errorf("unsupported number of channels: %d", cImage.channels)
	}

	return floatimage.NewFloat64NRGBA(image.Rect(0, 0, width, height), data), nil
}

// isHDR returns true if the file extension indicates an HDR format
func isHDR(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".hdr", ".exr", ".pfm", ".dpx":
		return true
	default:
		return false
	}
}

func WriteImage(filename string, image image.Image) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var cError *C.char
	defer func() {
		if cError != nil {
			C.free(unsafe.Pointer(cError))
		}
	}()

	isHDRFormat := isHDR(filename)
	cImage := toCImage(image, isHDRFormat)
	defer C.free_image(cImage)

	var cHdr C.int
	if isHDRFormat {
		cHdr = 1
	}

	if ret := C.write_image(cFilename, cImage, &cError, cHdr); ret != 0 {
		if cError != nil {
			return errors.New(C.GoString(cError))
		}
		return errors.New("failed to write image")
	}

	return nil
}

func toCImage(image image.Image, isHDRFormat bool) *C.Image {
	bounds := image.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	numChannels := 4

	// Allocate Image struct on the heap
	cImage := (*C.Image)(C.malloc(C.size_t(unsafe.Sizeof(C.Image{}))))
	cImage.width = C.int(width)
	cImage.height = C.int(height)
	cImage.channels = C.int(numChannels)
	cImage.data = (*C.float)(C.malloc(C.size_t(width) * C.size_t(height) * C.size_t(numChannels) * C.size_t(unsafe.Sizeof(float32(0)))))

	// Convert to a slice for easier access
	data := (*[1 << 30]C.float)(unsafe.Pointer(cImage.data))[: width*height*numChannels : width*height*numChannels]

	switch image := image.(type) {
	case *floatimage.Float32NRGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := image.Float32NRGBAAt(x, y)
				idx := ((y-bounds.Min.Y)*width + (x - bounds.Min.X)) * numChannels
				if isHDRFormat {
					// For HDR formats, pass the values directly
					data[idx] = C.float(c.R)
					data[idx+1] = C.float(c.G)
					data[idx+2] = C.float(c.B)
					data[idx+3] = C.float(c.A)
				} else {
					// For LDR formats, apply gamma correction and scale to [0,1]
					gamma := 2.2
					data[idx] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, float64(c.R))), 1.0/gamma))
					data[idx+1] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, float64(c.G))), 1.0/gamma))
					data[idx+2] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, float64(c.B))), 1.0/gamma))
					data[idx+3] = C.float(math.Min(1.0, math.Max(0.0, float64(c.A))))
				}
			}
		}
	case *floatimage.Float64NRGBA:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := image.Float64NRGBAAt(x, y)
				idx := ((y-bounds.Min.Y)*width + (x - bounds.Min.X)) * numChannels
				if isHDRFormat {
					// For HDR formats, pass the values directly
					data[idx] = C.float(c.R)
					data[idx+1] = C.float(c.G)
					data[idx+2] = C.float(c.B)
					data[idx+3] = C.float(c.A)
				} else {
					// For LDR formats, apply gamma correction and scale to [0,1]
					gamma := 2.2
					data[idx] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, c.R)), 1.0/gamma))
					data[idx+1] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, c.G)), 1.0/gamma))
					data[idx+2] = C.float(math.Pow(math.Min(1.0, math.Max(0.0, c.B)), 1.0/gamma))
					data[idx+3] = C.float(math.Min(1.0, math.Max(0.0, c.A)))
				}
			}
		}
	default:
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := image.At(x, y).RGBA()
				idx := ((y-bounds.Min.Y)*width + (x - bounds.Min.X)) * numChannels
				// Convert from 16-bit RGBA to float in range [0,1]
				data[idx] = C.float(float32(r) / 65535.0)
				data[idx+1] = C.float(float32(g) / 65535.0)
				data[idx+2] = C.float(float32(b) / 65535.0)
				data[idx+3] = C.float(float32(a) / 65535.0)
			}
		}
	}

	return cImage
}

func toRGBA32Slice(cData []C.float, width int, height int, numChannels int) []float32 {

	data := make([]float32, width*height*numChannels)
	for i := 0; i < len(cData); i++ {
		data[i] = float32(cData[i])
	}
	return data
}

func toRGB32Slice(cData []C.float, width int, height int) []float32 {
	data := make([]float32, width*height*4)

	j := 0

	for i := 0; i < len(data); i += 4 {
		data[i] = float32(cData[j])
		j++
		data[i+1] = float32(cData[j])
		j++
		data[i+2] = float32(cData[j])
		j++
		data[i+3] = 1.0 // alpha channel is always 1.0
	}

	return data
}

func toRGBA64Slice(cData []C.float, width int, height int, numChannels int) []float64 {

	data := make([]float64, width*height*numChannels)
	for i := 0; i < len(cData); i++ {
		data[i] = float64(cData[i])
	}
	return data
}

func toRGB64Slice(cData []C.float, width int, height int) []float64 {
	data := make([]float64, width*height*4)

	j := 0

	for i := 0; i < len(data); i += 4 {
		data[i] = float64(cData[j])
		j++
		data[i+1] = float64(cData[j])
		j++
		data[i+2] = float64(cData[j])
		j++
		data[i+3] = 1.0 // alpha channel is always 1.0
	}

	return data
}
