package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math/cmplx"
	"os"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	windowName                  = "explorer"
	maxTicksPerSec              = 30
	windowWidth, windowHeight   = 800, 600
	fWindowWidth, fWindowHeight = float64(windowWidth), float64(windowHeight)
	fWindowWidthDiv2            = fWindowWidth / 2.0
	fWindowHeightDiv2           = fWindowHeight / 2.0
)

var (
	maxIterations      int = 300
	fMaxIterations         = float64(maxIterations)
	xpos, ypos         float64
	zoom               float64 = 1
	renderJulia        bool
	fastEvalEnabled    bool
	beginViz           bool
	iterationBuffer    []int
	frameBuffer        []byte
	colors             []color.RGBA
	lastRenderDuration time.Duration
	mplusNormalFont    font.Face
)

func init() {
	parseFlags()

	iterationBuffer = make([]int, windowWidth*windowHeight)
	// Need 4 bytes (r,g,b,a) for each pixel which is colored per frame.
	frameBuffer = make([]byte, windowWidth*windowHeight*4)
	colors = interpolateColors(fMaxIterations)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Game is the required type from ebiten which must implement that package's
// expected game loop functions.
type Game struct{}

// Update is called on every loop "tick". Ebiten will attempt to call this up to
// the max allowable TPS, but due to the high cost of our rendering function
// ticks per second will end up being much less than the 60/sec default.
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if !beginViz {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			beginViz = true
		}
	}

	shiftAmt := 0.1 / zoom

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		ypos += shiftAmt
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		ypos -= shiftAmt
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		xpos += shiftAmt
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		xpos -= shiftAmt
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		xpos, ypos = 0.0, 0.0
		zoom = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyO) {
		zoom -= zoom * 0.03
	} else if ebiten.IsKeyPressed(ebiten.KeyI) {
		zoom += zoom * 0.03
	}
	if zoom == 0 {
		zoom = 1
	}
	// zoom += yScrollOffset

	renderFrame()
	return nil
}

// Draw is called on every frame and updates the ebiten screen image.
func (g *Game) Draw(screen *ebiten.Image) {
	if !beginViz {
		text.Draw(
			screen,
			"Move with arrow keys\n"+
				"Zoom in/out with 'I'/'O' keys\n"+
				"Reset with 'R' key\n"+
				"Exit with 'Escape'\n"+
				"Press Enter to start\n",
			mplusNormalFont,
			20, 40, color.White,
		)
		return
	}

	screen.ReplacePixels(frameBuffer)
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"FPS: %0.2f\nTPS: %0.2f\nLast Render Time: %v\nZoom: %0.2f\nXpos: %0.2f\nYpos: %0.2f",
			ebiten.CurrentFPS(),
			ebiten.CurrentTPS(),
			lastRenderDuration,
			zoom, xpos, ypos,
		),
	)
}

// Layout changes the screen size based on users changing the window size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func julia(z complex128) int {
	c := complex(0.25, 0.5)
	var i int
	for i = 0; i < maxIterations; i++ {
		if cmplx.Abs(z) > 2 {
			return i
		}
		z = cmplx.Pow(z, 2) + c
	}

	return maxIterations
}

func mandel(z complex128) int {
	newZ := z
	var i int
	for i = 0; i < maxIterations; i++ {
		if cmplx.Abs(newZ) > 2 {
			return i
		}
		newZ = (newZ * newZ) + z
	}

	return maxIterations
}

func parseFlags() {
	flag.Float64Var(
		&ypos, "y position", 0,
		"Starting position along the imaginary axis of the complex plane")
	flag.Float64Var(
		&xpos, "x position", 0,
		"Starting position along the real axis of the complex plane")
	flag.BoolVar(
		&renderJulia, "julia", false,
		"Visualize a Julia set, this is really slow, don't use it")
	flag.BoolVar(
		&fastEvalEnabled, "fast eval", true,
		"Use an evaluation estimation for a render speedup")
	flag.Parse()
}

// useNeighborFastEval uses historical iteration counts for each pixel to determine
// if a particular pixel can be colored the same as all of it's neighbors. As long
// as all a pixels neighbors have the same iteration result, then that result is
// returned. Otherwise a nil error is returned.
//
// Use of this function allows for optimizing frame updates in exchange for lower
// resolution rendering as the user moves.
func useNeighborFastEval(x, y int) (int, error) {
	if !fastEvalEnabled {
		return 0, errors.New("Fast eval disabled")
	}

	left := (x - 1) + (y * windowWidth)
	right := (x + 1) + (y * windowWidth)
	up := x + ((y + 1) * windowWidth)
	down := x + ((y - 1) * windowWidth)
	if x > 0 && x < windowWidth-1 && y > 0 && y < windowHeight-1 {
		if iterationBuffer[left] == iterationBuffer[right] &&
			iterationBuffer[up] == iterationBuffer[down] &&
			iterationBuffer[left] == iterationBuffer[up] {
			return iterationBuffer[left], nil
		}
	}
	return 0, errors.New("Can't use neighbors")
}

// renderFrame draws one frame of the image to frameBuffer, checking each pixel at
// the current location and zoom level to see if it is bounded or not.
func renderFrame() {
	// Each row of the output is computed in parallel goroutines.
	wg := sync.WaitGroup{}
	start := time.Now()
	for y := 0; y < windowHeight; y++ {
		// Scale y from (0, windowHeight) to the plane size, depending on zoom and ypos.
		yi := ((float64(y) - fWindowHeightDiv2) /
			(0.5 * zoom * fWindowHeight)) + ypos
		wg.Add(1)
		go func(yi float64, y int) {
			defer wg.Done()
			for x := 0; x < windowWidth; x++ {
				var iterCount int

				iterations, err := useNeighborFastEval(x, y)
				if err == nil {
					iterCount = iterations
				} else {
					// Scale y from (0, windowHeight) to the plane size,
					// depending on zoom and ypos.
					xi := (1.5 * (float64(x) - fWindowWidthDiv2) /
						(0.5 * zoom * fWindowWidth)) + xpos
					z := complex(xi, yi)
					if renderJulia {
						iterCount = julia(z)
					} else {
						iterCount = mandel(z)
					}

				}

				// Cache result for fastEval checks.
				// This is an intential datarace! Locking this slice slows
				// down rendering too much because multiple goroutines need
				// to then synchronize in order to finish this write. This
				// hasn't crashed yet, but it probably should.
				iterationBuffer[x+(y*windowWidth)] = iterCount

				// Use black for high iteration counts.
				pixelColor := color.RGBA{}
				if iterCount < len(colors)-1 {
					color1 := colors[iterCount]
					color2 := colors[iterCount+1]
					col := linearInterpolation(
						rgbaToUint(color1), rgbaToUint(color2), uint32(iterCount))
					pixelColor = uint32ToRgba(col)

				}
				p := 4 * (x + (y * windowWidth))
				frameBuffer[p] = pixelColor.R
				frameBuffer[p+1] = pixelColor.G
				frameBuffer[p+2] = pixelColor.B
				frameBuffer[p+3] = 0xff // Alpha is always 255
			}
		}(yi, y)
	}
	wg.Wait()
	lastRenderDuration = time.Since(start)
}

// main creates an ebiten game window and begins the game loop.
func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(windowName)
	ebiten.SetMaxTPS(maxTicksPerSec)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
