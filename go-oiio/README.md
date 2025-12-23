# Golang bindings for OpenImageIO


## Introduction
This package provides Go bindings to read and write images using [OpenImageIO](https://openimageio.readthedocs.io).

The image data is always in [floatimage](https://github.com/flynn-nrg/floatimage) format as the use case when building this was to allow [Izpi](https://github.com/flynn-nrg/izpi) to easily import texture assets in many different formats.

## Building
Use your package manager to install the following dependencies:

* pkgconf
* fmt
* OpenImageIO

You can verify this by running the following command:

```shell
$ pkgconf --list-all| egrep '(OpenImage|fmt)'
fmt                            fmt - A modern formatting library
OpenImageIO                    OpenImageIO - OpenImageIO is a library for reading and writing images.
```

## Example usage

### Basic Image I/O

```golang
package main

import (
	"fmt"
	"log"

	"github.com/flynn-nrg/go-vfx/go-oiio/oiio"
)

func main() {
	floatImage32, err := oiio.ReadImage32("test.exr")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("image size: %d x %d\n", floatImage32.Bounds().Dx(), floatImage32.Bounds().Dy())

	if err := oiio.WriteImage("test.png", floatImage32); err != nil {
		log.Fatal(err)
	}
}
```

### Color Space Conversions with OpenColorIO

The package supports automatic color space conversions using OpenColorIO (OCIO) integration.

**Requirements:**
- OpenColorIO 2.x installed
- Either use the built-in config (`export OCIO=ocio://default`) or set `OCIO` to your config file path

#### Available Conversion Options

- **`oiio.Raw`** - Read image data as-is without any conversion
- **`oiio.LineariseSRGB`** - Convert sRGB to linear Rec.709 (removes gamma, keeps primaries)
- **`oiio.ConvertToACEScg`** - Convert sRGB to ACEScg (AP1) color space (linearises and converts primaries)

#### Reading Images with Color Space Conversion

```golang
package main

import (
	"fmt"
	"log"

	"github.com/flynn-nrg/go-vfx/go-oiio/oiio"
)

func main() {
	// Option 1: Read raw image data without conversion
	imgRaw, err := oiio.ReadImage32("texture.png")
	if err != nil {
		log.Fatal(err)
	}
	
	// Option 2: Read and linearize sRGB to linear Rec.709
	imgLinear, err := oiio.ReadImage("texture.png", oiio.LineariseSRGB)
	if err != nil {
		log.Fatal(err)
	}
	
	// Option 3: Read and convert to ACEScg (AP1) color space
	// This automatically:
	// 1. Linearizes sRGB-encoded data (applies inverse gamma)
	// 2. Converts from Rec. 709 primaries to ACEScg (AP1) primaries
	imgAces, err := oiio.ReadImage("texture.png", oiio.ConvertToACEScg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ACEScg image size: %d x %d\n", imgAces.Bounds().Dx(), imgAces.Bounds().Dy())
	
	// The image data is now in linear ACEScg color space
	// suitable for rendering pipelines using ACES workflow
}
```

### Complete ACEScg Workflow Example

```golang
package main

import (
	"image"
	"log"

	"github.com/flynn-nrg/go-vfx/go-oiio/oiio"
)

func main() {
	// Read image and convert to ACEScg color space
	imgAces, err := oiio.ReadImage("input.png", oiio.ConvertToACEScg)
	if err != nil {
		log.Fatal(err)
	}

	// Process the image in ACEScg color space
	// (your rendering/compositing code here)

	// Prepare ACES metadata for output
	bounds := imgAces.Bounds()
	metadata := &oiio.ACESMetadata{
		DisplayWindow:    image.Rect(0, 0, 1920, 1080), // Full canvas
		DataWindow:       bounds,                        // Actual pixel data region
		PixelAspectRatio: 1.0,                          // 1.0 for square pixels
		Timecode:         "01:23:45:12",                // Optional SMPTE timecode
		ACESVersion:      "ACES 1.3",                   // ACES version
	}

	// Write as EXR with full ACES metadata
	if err := oiio.WriteImageACES("output.exr", imgAces, metadata); err != nil {
		log.Fatal(err)
	}
}
```

### ACEScg Output Features

When using `WriteImageACES`, the following metadata is embedded in the EXR file:

- **Color Space**: `oiio:ColorSpace` attribute set to "ACEScg"
- **Chromaticities**: ACEScg (AP1) color primaries and D60 white point:
  - Red primary: (0.713, 0.300)
  - Green primary: (0.165, 0.830)
  - Blue primary: (0.128, 0.044)
  - White point: (0.32168, 0.33767) - D60
- **Display/Data Windows**: Support for overscan, tiling, and cropping
- **Pixel Aspect Ratio**: Configurable (1.0 for square pixels, 2.0 for 2x anamorphic, etc.)
- **Timecode**: Optional SMPTE timecode for frame identification
- **ACES Version**: Custom attribute specifying ACES version
- **Compression**: ZIP compression (lossless, industry standard)

## Limitations

* Only 3 and 4 channel images are supported for now.
* [OpenEXR](https://openexr.com) deep pixels are not yet supported.