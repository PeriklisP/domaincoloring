package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"

	"github.com/fogleman/colormap"
	"github.com/nfnt/resize"
)

const (
	W = 1600
	H = 1600

	CenterReal = 0
	CenterImag = 0

	Fovy = 1

	Supersampling = 2

	Sw = Supersampling * W
	Sh = Supersampling * H

	AspectRatio = float64(W) / H
	HalfFovy    = float64(Fovy) / 2
)

var (
	Colormap = colormap.Inferno
)

func complexFunction(z complex128) complex128 {
	return cmplx.Sin(1 / z)
}

func complexColor(z complex128) color.Color {
	modulus := cmplx.Abs(z)
	phase := cmplx.Phase(z)

	t0 := modulus / (modulus + 1)

	// t1 := math.Sin(phase)
	// t1 = (t1 + 1) / 2

	t1 := phase/math.Pi + 1
	if t1 > 1 {
		t1 = 2 - t1
	}

	t := math.Pow(t0, 2) * t1
	// t = t0

	// t = ease.InOutQuad(t)

	return Colormap.At(t)
}

func main() {
	im := image.NewNRGBA64(image.Rect(0, 0, Sw, Sh))

	for py := 0; py < Sh; py++ {
		y := ((float64(Sh-py-1)/(Sh-1))*2-1)*HalfFovy + CenterImag
		for px := 0; px < Sw; px++ {
			x := ((float64(px)/(Sw-1))*2-1)*AspectRatio*HalfFovy + CenterReal
			z := complexFunction(complex(x, y))
			c := complexColor(z)
			im.Set(px, py, c)
		}
	}

	downsampled := resize.Resize(W, H, im, resize.Bilinear)

	err := savePNG("out.png", downsampled)
	if err != nil {
		log.Fatal(err)
	}
}

func savePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}
