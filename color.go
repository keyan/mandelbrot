// This code is copied from https://github.com/esimov/gobrot because I really liked
// his color interpolation scheme, but couldn't figure out exactly how the interpolation
// calculation worked.
package main

import (
	"image/color"
	"math"
)

var colorPalette = []color.RGBA{
	{0x00, 0x04, 0x0f, 0xff},
	{0x03, 0x26, 0x28, 0xff},
	{0x07, 0x3e, 0x1e, 0xff},
	{0x18, 0x55, 0x08, 0xff},
	{0x5f, 0x6e, 0x0f, 0xff},
	{0x84, 0x50, 0x19, 0xff},
	{0x9b, 0x30, 0x22, 0xff},
	{0xb4, 0x92, 0x2f, 0xff},
	{0x94, 0xca, 0x3d, 0xff},
	{0x4f, 0xd5, 0x51, 0xff},
	{0x66, 0xff, 0xb3, 0xff},
	{0x82, 0xc9, 0xe5, 0xff},
	{0x9d, 0xa3, 0xeb, 0xff},
	{0xd7, 0xb5, 0xf3, 0xff},
	{0xfd, 0xd6, 0xf6, 0xff},
	{0xff, 0xf0, 0xf2, 0xff},
}

func interpolateColors(numberOfColors float64) []color.RGBA {
	var factor float64
	steps := []float64{}
	cols := []uint32{}
	interpolated := []uint32{}
	interpolatedColors := []color.RGBA{}

	factor = 1.0 / numberOfColors
	for index, col := range colorPalette {
		stepRatio := float64(index+1) / float64(len(colorPalette))
		step := float64(int(stepRatio*100)) / 100
		steps = append(steps, step)
		r, g, b, a := col.RGBA()
		r /= 0xff
		g /= 0xff
		b /= 0xff
		a /= 0xff
		uintColor := uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
		cols = append(cols, uintColor)
	}

	var min, max, minColor, maxColor float64
	if len(colorPalette) == len(steps) && len(colorPalette) == len(cols) {
		for i := 0.0; i <= 1; i += factor {
			for j := 0; j < len(colorPalette)-1; j++ {
				if i >= steps[j] && i < steps[j+1] {
					min = steps[j]
					max = steps[j+1]
					minColor = float64(cols[j])
					maxColor = float64(cols[j+1])
					uintColor := cosineInterpolation(maxColor, minColor, (i-min)/(max-min))
					interpolated = append(interpolated, uint32(uintColor))
				}
			}
		}
	}

	for _, pixelValue := range interpolated {
		r := pixelValue >> 24 & 0xff
		g := pixelValue >> 16 & 0xff
		b := pixelValue >> 8 & 0xff
		a := 0xff

		interpolatedColors = append(interpolatedColors, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	return interpolatedColors
}

func cosineInterpolation(c1, c2, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2.0
	return c1*(1-mu2) + c2*mu2
}

func linearInterpolation(c1, c2, mu uint32) uint32 {
	return c1*(1-mu) + c2*mu
}

func rgbaToUint(color color.RGBA) uint32 {
	r, g, b, a := color.RGBA()
	r /= 0xff
	g /= 0xff
	b /= 0xff
	a /= 0xff
	return uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
}

func uint32ToRgba(col uint32) color.RGBA {
	r := col >> 24 & 0xff
	g := col >> 16 & 0xff
	b := col >> 8 & 0xff
	a := 0xff
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
