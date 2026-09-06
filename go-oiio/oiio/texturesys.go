package oiio

/*
#cgo CXXFLAGS: -std=c++17
#cgo pkg-config: OpenImageIO fmt
#include <stdlib.h>
#include "./texturesys_wrapper.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// WrapMode selects how texture lookups behave outside the [0,1] st range.
type WrapMode C.TextureWrapMode

const (
	WrapDefault  WrapMode = C.WRAP_DEFAULT
	WrapBlack    WrapMode = C.WRAP_BLACK
	WrapClamp    WrapMode = C.WRAP_CLAMP
	WrapPeriodic WrapMode = C.WRAP_PERIODIC
	WrapMirror   WrapMode = C.WRAP_MIRROR
)

// TextureLookupOptions mirrors OIIO's TextureOpt fields relevant to a single
// filtered 2D lookup. Use NewTextureLookupOptions to get OIIO's defaults
// (notably SWidth/TWidth of 1) rather than a zero-valued struct.
type TextureLookupOptions struct {
	FirstChannel   int
	Subimage       int
	SWrap, TWrap   WrapMode
	SBlur, TBlur   float32
	SWidth, TWidth float32
	Fill           float32
}

// NewTextureLookupOptions returns TextureLookupOptions set to OIIO's own
// defaults. Starting from a zero-valued struct instead would set SWidth/TWidth
// to 0, disabling derivative-based filtering rather than using OIIO's default
// multiplier of 1.
func NewTextureLookupOptions() TextureLookupOptions {
	return TextureLookupOptions{SWidth: 1, TWidth: 1}
}

func (o TextureLookupOptions) toC() C.TextureLookupOptions {
	return C.TextureLookupOptions{
		firstchannel: C.int(o.FirstChannel),
		subimage:     C.int(o.Subimage),
		swrap:        C.TextureWrapMode(o.SWrap),
		twrap:        C.TextureWrapMode(o.TWrap),
		sblur:        C.float(o.SBlur),
		tblur:        C.float(o.TBlur),
		swidth:       C.float(o.SWidth),
		twidth:       C.float(o.TWidth),
		fill:         C.float(o.Fill),
	}
}

// TextureSystem wraps an OIIO::TextureSystem: a filtered, mipmapped,
// tile-cached texture lookup engine, as opposed to the whole-image I/O in
// ReadImage/WriteImage. Callers must call Close when done; a finalizer is
// registered as a safety net but should not be relied upon since the
// underlying tile cache can be large.
type TextureSystem struct {
	ptr *C.TextureSystemHandle
}

// NewTextureSystem creates a TextureSystem backed by OIIO's shared tile
// cache. It is safe to create more than one; OIIO shares the underlying
// cache across instances created with shared=true.
func NewTextureSystem() (*TextureSystem, error) {
	var errorMsg *C.char
	ptr := C.texturesystem_create(&errorMsg)
	if ptr == nil {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return nil, fmt.Errorf("failed to create texture system: %s", err)
		}
		return nil, fmt.Errorf("failed to create texture system")
	}

	ts := &TextureSystem{ptr: ptr}
	runtime.SetFinalizer(ts, (*TextureSystem).Close)
	return ts, nil
}

// Close releases the underlying OIIO::TextureSystem and its tile cache. It
// is safe to call more than once.
func (ts *TextureSystem) Close() {
	if ts.ptr != nil {
		C.texturesystem_destroy(ts.ptr)
		ts.ptr = nil
		runtime.SetFinalizer(ts, nil)
	}
}

// Texture performs a filtered lookup at (s,t), using the screen-space
// derivatives dsdx/dtdx (d/dx) and dsdy/dtdy (d/dy) to pick a mip level and
// filter width. This mirrors OSL's texture() shadeop, which always supplies
// derivatives computed from the shading point's surface differentials.
// Returns nchannels float32 values.
func (ts *TextureSystem) Texture(filename string, opts TextureLookupOptions, s, t, dsdx, dtdx, dsdy, dtdy float32, nchannels int) ([]float32, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cOpts := opts.toC()
	result := make([]float32, nchannels)

	var errorMsg *C.char
	ret := C.texturesystem_texture(ts.ptr, cFilename, &cOpts, C.float(s), C.float(t),
		C.float(dsdx), C.float(dtdx), C.float(dsdy), C.float(dtdy),
		C.int(nchannels), (*C.float)(unsafe.Pointer(&result[0])), &errorMsg)
	if ret != 0 {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return nil, fmt.Errorf("texture lookup failed: %s", err)
		}
		return nil, fmt.Errorf("texture lookup failed")
	}

	return result, nil
}

// GetTextureInfoInt retrieves an integer-valued texture attribute, e.g.
// "channels" or "exists", mirroring OSL's gettextureinfo() shadeop.
func (ts *TextureSystem) GetTextureInfoInt(filename, dataname string) (int, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	cDataname := C.CString(dataname)
	defer C.free(unsafe.Pointer(cDataname))

	var out C.int
	var errorMsg *C.char
	ret := C.texturesystem_get_texture_info_int(ts.ptr, cFilename, cDataname, &out, &errorMsg)
	if ret != 0 {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return 0, fmt.Errorf("get_texture_info failed: %s", err)
		}
		return 0, fmt.Errorf("get_texture_info failed")
	}

	return int(out), nil
}

// GetTextureInfoFloat retrieves a float-valued texture attribute, mirroring
// OSL's gettextureinfo() shadeop.
func (ts *TextureSystem) GetTextureInfoFloat(filename, dataname string) (float32, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	cDataname := C.CString(dataname)
	defer C.free(unsafe.Pointer(cDataname))

	var out C.float
	var errorMsg *C.char
	ret := C.texturesystem_get_texture_info_float(ts.ptr, cFilename, cDataname, &out, &errorMsg)
	if ret != 0 {
		if errorMsg != nil {
			err := C.GoString(errorMsg)
			C.free(unsafe.Pointer(errorMsg))
			return 0, fmt.Errorf("get_texture_info failed: %s", err)
		}
		return 0, fmt.Errorf("get_texture_info failed")
	}

	return float32(out), nil
}
