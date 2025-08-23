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

## Limitations

* Only 3 and 4 channel images are supported for now.
* [OpenEXR](https://openexr.com) deep pixels are not yet supported.