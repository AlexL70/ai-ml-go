package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/StephaneBunel/bresenham"
	"github.com/kmicki/apng"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Constant
const cellSize = 60

// Variables for color
var (
	green     = color.RGBA{G: 255, A: 255}
	darkGreen = color.RGBA{G: 100, R: 1, B: 32, A: 255}
	red       = color.RGBA{R: 255, A: 255}
	yellow    = color.RGBA{R: 255, G: 255, B: 101, A: 255}
	gray      = color.RGBA{R: 125, G: 125, B: 125, A: 255}
	orange    = color.RGBA{R: 255, G: 140, B: 25, A: 255}
	blue      = color.RGBA{R: 40, G: 118, B: 173, A: 255}
)

// OutputImage — draw the maze as png file.
func (g *Maze) OutputImage(fileName ...string) {
	fmt.Printf("Generating image %s...\n", fileName)
	width := cellSize * g.Width
	height := cellSize * g.Height
	var outFile = "image.png"
	if len(fileName) > 0 {
		outFile = fileName[0]
	}

	uplLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: uplLeft, Max: lowRight})
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
	// draw squares on the image
	for i, row := range g.Walls {
		for j, col := range row {
			p := Point{Row: i, Col: j}
			if col.Wall {
				// black square for wall
				g.drawSquare(col, p, img, color.Black, cellSize, j*cellSize, i*cellSize)
			} else if g.InSolution(p) {
				// green square for the part of solution
				g.drawSquare(col, p, img, green, cellSize, j*cellSize, i*cellSize)
			} else if p.Col == g.Start.Col && p.Row == g.Start.Row {
				// dark green square for the starting point
				g.drawSquare(col, p, img, darkGreen, cellSize, j*cellSize, i*cellSize)
			} else if p.Col == g.Goal.Col && p.Row == g.Goal.Row {
				// red for the goal
				g.drawSquare(col, p, img, red, cellSize, j*cellSize, i*cellSize)
			} else if col.State == g.CurrentNode.State {
				// orange for the current location
				g.drawSquare(col, p, img, orange, cellSize, j*cellSize, i*cellSize)
			} else if inExplored(p, g.Explored) {
				// Yellow for explored locations
				g.drawSquare(col, p, img, yellow, cellSize, j*cellSize, i*cellSize)
			} else {
				// white for empty unexplored
				g.drawSquare(col, p, img, color.White, cellSize, j*cellSize, i*cellSize)
			}
		}
	}

	// Draw a grid
	for i, _ := range g.Walls {
		bresenham.DrawLine(img, 0, i*cellSize, g.Width*cellSize, i*cellSize, gray)
	}
	for i := 0; i <= g.Width; i++ {
		bresenham.DrawLine(img, i*cellSize, 0, i*cellSize, g.Height*cellSize, gray)
	}

	f, err := os.Create(outFile)
	if err != nil {
		panic(fmt.Errorf("error creating image outfile: %w", err))
	}
	defer f.Close()
	_ = png.Encode(f, img)
}

// drawSquare
func (g *Maze) drawSquare(col Wall, p Point, img *image.RGBA, c color.Color, size, x, y int) {
	patch := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(patch, patch.Bounds(), &image.Uniform{
		C: c,
	}, image.Point{}, draw.Src)

	if !col.Wall {
		g.printLocation(p, color.Black, patch)
	}

	draw.Draw(img, image.Rect(x, y, x+size, y+size), patch, image.Point{}, draw.Src)

}

// printLocation
func (g *Maze) printLocation(p Point, c color.Color, patch *image.RGBA) {
	point := fixed.Point26_6{X: fixed.I(6), Y: fixed.I(40)}
	d := &font.Drawer{
		Dst:  patch,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(fmt.Sprintf("[%d %d]", p.Row, p.Col))
}

func (g *Maze) OutputAnimatedImage() {
	output := "./animation.png"
	files, err := os.ReadDir("./tmp")
	if err != nil {
		panic(fmt.Errorf("error opening animation sources: %w", err))
	}
	var images []string
	var delays []int
	for _, file := range files {
		images = append(images, fmt.Sprintf("./tmp/%s", file.Name()))
		delays = append(delays, 30)
	}
	images = append(images, "./image.png")
	a := apng.APNG{
		Frames: make([]apng.Frame, len(images)),
	}
	out, err := os.Create(output)
	if err != nil {
		panic(fmt.Errorf("error creating animation file: %w", err))
	}
	defer out.Close()

	for i, s := range images {
		in, err := os.Open(s)
		if err != nil {
			panic(fmt.Errorf("error reading images: %w", err))
		}
		defer in.Close()

		m, err := png.Decode(in)
		if err != nil {
			// it is continue in the original code in case if arbitrary file appear
			panic(fmt.Errorf("error decoding image: %w", err))
		}
		a.Frames[i].Image = m
	}

	err = apng.Encode(out, a)
	if err != nil {
		panic(fmt.Errorf("error encoding images: %w", err))
	}
}
