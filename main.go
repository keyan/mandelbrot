package main

import (
	"flag"
	"fmt"
	"runtime"

	"math/cmplx"

	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"

	"gopkg.in/teh-cmc/go-sfml.v24/graphics"
	"gopkg.in/teh-cmc/go-sfml.v24/window"
)

const (
	windowName    = "explorer"
	maxIterations = 300
)

var (
	windowWidth, windowHeight int = 800, 600
	xpos, ypos                float64
	zoom                      float64 = 1
	renderJulia               bool
)

// julia returns a value <= 1 representing the number of iterations
// (relative to the max) it took for a complex number to escape to infinity.
// If a number remains bounded then 1.0 is returned.
func julia(z complex128) float64 {
	c := complex(0.25, 0.5)
	i := 0
	for ; i < maxIterations; i++ {
		if cmplx.Abs(z) > 2 {
			break
		}
		z = cmplx.Pow(z, 2) + c
	}

	return float64(maxIterations-i) / maxIterations
}

// mandel returns a value <= 1 representing the number of iterations
// (relative to the max) it took for a complex number to escape to infinity.
// If a number remains bounded then 1.0 is returned.
func mandel(z complex128) float64 {
	newZ := z
	i := 0
	for ; i < maxIterations; i++ {
		if cmplx.Abs(newZ) > 2 {
			break
		}
		newZ = (newZ * newZ) + z
	}

	return float64(maxIterations-i) / maxIterations
}

func parseFlags() {
	flag.IntVar(&windowWidth, "windowWidth", windowWidth, "Width of output image")
	flag.IntVar(&windowHeight, "windowHeight", windowHeight, "Height of output image")
	flag.Float64Var(
		&ypos, "y position", 0,
		"Starting position along the imaginary axis of the complex plane")
	flag.Float64Var(
		&xpos, "x position", 0,
		"Starting position along the real axis of the complex plane")
	flag.BoolVar(
		&renderJulia, "julia", false, "Visualize a Julia set")
	flag.Parse()
}

func renderFrame() {
	wg := sync.WaitGroup{}
	img := image.NewNRGBA(image.Rect(0, 0, windowWidth, windowHeight))
	start := time.Now()
	for y := 0; y <= windowHeight; y++ {
		yi := (float64(y) - (float64(windowHeight) / 2)) / ((0.5 * zoom * float64(windowHeight)) + ypos)
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x <= windowWidth; x++ {
				xi := 1.5 * (float64(x) - (float64(windowWidth) / 2)) / ((0.5 * zoom * float64(windowWidth)) + xpos)
				z := complex(xi, yi)
				var escapeVal float64
				if renderJulia {
					escapeVal = julia(z)
				} else {
					escapeVal = mandel(z)
				}
				img.Set(x, y, color.NRGBA{
					R: uint8(escapeVal * 230),
					G: uint8(escapeVal * 235),
					B: uint8(escapeVal * 255),
					A: 255,
				})
			}
		}(y)
	}
	wg.Wait()
	fmt.Printf("Render time: %v\n", time.Since(start))

	f, err := os.Create("output.png")
	defer f.Close()
	if err != nil {
		fmt.Println("Failed to create output file")
		os.Exit(1)
	}

	if err := png.Encode(f, img); err != nil {
		fmt.Println("Failed to encode image")
		os.Exit(1)
	}
}

func init() {
	// Required for SFML
	runtime.LockOSThread()
}

func main() {
	parseFlags()
	renderFrame()

	// Window sizing
	vm := window.NewSfVideoMode()
	defer window.DeleteSfVideoMode(vm)
	vm.SetWidth(uint(windowWidth))
	vm.SetHeight(uint(windowHeight))
	vm.SetBitsPerPixel(32)

	// Main window
	cs := window.NewSfContextSettings()
	defer window.DeleteSfContextSettings(cs)
	w := graphics.SfRenderWindow_create(vm, windowName, uint(window.SfResize|window.SfClose), cs)
	defer window.SfWindow_destroy(w)

	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	// Start the game loop
	for window.SfWindow_isOpen(w) > 0 {
		// Process events
		for window.SfWindow_pollEvent(w, ev) > 0 {
			// Close window: exit
			if ev.GetXtype() == window.SfEventType(window.SfEvtClosed) {
				return
			}
		}

		renderFrame()
		//https://www.sfml-dev.org/documentation/2.3.2/classsf_1_1Image.php#a1c2b960ea12bdbb29e80934ce5268ebf
		// graphics.SfImageCreate()
		graphics.SFImage_createFromPixels(
		graphics.SfRenderWindow_clear(w, graphics.GetSfRed())
		graphics.SfRenderWindow_display(w)
	}
}
