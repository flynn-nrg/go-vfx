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

### ACEScg Color Space Workflow

The package includes full support for ACEScg (AP1) color space workflow, which is the industry standard for VFX and animation pipelines:

#### Reading Images in ACEScg

```golang
package main

import (
	"fmt"
	"log"

	"github.com/flynn-nrg/go-vfx/go-oiio/oiio"
)

func main() {
	// Read image and convert to ACEScg (AP1) color space
	// This automatically:
	// 1. Linearizes sRGB-encoded data (applies inverse gamma)
	// 2. Converts from Rec. 709 primaries to ACEScg (AP1) primaries
	imgAces, err := oiio.ReadImage32Aces("texture.png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ACEScg image size: %d x %d\n", imgAces.Bounds().Dx(), imgAces.Bounds().Dy())
	
	// The image data is now in linear ACEScg color space
	// suitable for rendering pipelines using ACES workflow
}
```

#### Writing Images in ACEScg with Metadata

```golang
package main

import (
	"image"
	"log"

	"github.com/flynn-nrg/go-vfx/go-oiio/oiio"
)

func main() {
	// Read image in ACEScg color space
	imgAces, err := oiio.ReadImage32Aces("input.png")
	if err != nil {
		log.Fatal(err)
	}

	// Process the image...
	// (your rendering/compositing code here)

	// Prepare ACES metadata
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

## Limitations

* Only 3 and 4 channel images are supported for now.
* [OpenEXR](https://openexr.com) deep pixels are not yet supported.