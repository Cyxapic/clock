package main

/*
Стырил отсюдава:
	https://github.com/cnet-sudo/Clock/blob/clock1/main.cpp
Статья на хабре:
	https://habr.com/ru/post/706954/
*/

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCR_WIDTH  = 800
	SCR_HEIGHT = 800
	M_PI       = 3.14
)

//go:embed assets/clockBigArrow.png
var clockBigArrow []byte

//go:embed assets/clockBigMinArrow.png
var clockBigMinArrow []byte

//go:embed assets/clockBigHourArrow.png
var clockHourArrow []byte

//go:embed assets/bgClock.png
var bgClock []byte

//go:embed assets/bubbles.png
var bubbles []byte

//go:embed assets/smallPoint.png
var smallPoint []byte

func main() {
	ebiten.SetWindowSize(SCR_WIDTH, SCR_HEIGHT)
	ebiten.SetWindowTitle("Часики")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	BG        *BGClock
	secArrow  *Arrow
	minArrow  *Arrow
	hourArrow *Arrow
}

func NewGame() *Game {
	return &Game{
		NewBGClock(),
		NewArrow(clockBigArrow),
		NewArrow(clockBigMinArrow),
		NewArrow(clockHourArrow),
	}
}

func (g *Game) Update() error {
	hours, minutes, sec := time.Now().Clock()

	g.secArrow.RotateMultiplier = float64(sec) * 6
	g.minArrow.RotateMultiplier = float64(minutes)*6 + float64(sec)*0.1
	g.hourArrow.RotateMultiplier = 30*float64(hours) + float64(minutes)*0.5

	g.secArrow.Update()
	g.minArrow.Update()
	g.hourArrow.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.BG.Img, &ebiten.DrawImageOptions{})
	g.hourArrow.Draw(screen)
	g.minArrow.Draw(screen)
	g.secArrow.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCR_WIDTH, SCR_HEIGHT
}

type BGClock struct {
	Img *ebiten.Image
}

func NewBGClock() *BGClock {
	img := getImage(bgClock)
	drawClockPoints(img)
	return &BGClock{
		Img: img,
	}
}

type Arrow struct {
	Img              *ebiten.Image
	Op               *ebiten.DrawImageOptions
	RotateMultiplier float64
}

func NewArrow(arrow []byte) *Arrow {
	return &Arrow{
		Img: getImage(arrow),
		Op:  &ebiten.DrawImageOptions{},
	}
}

func (arrow *Arrow) Update() {
	w, h := arrow.Img.Bounds().Size().X, arrow.Img.Bounds().Size().Y
	arrow.Op = &ebiten.DrawImageOptions{}
	arrow.Op.GeoM.Translate(-float64(w)/2, -float64(h)+21)
	arrow.Op.GeoM.Rotate(arrow.RotateMultiplier * 2 * math.Pi / 360)
	arrow.Op.GeoM.Translate(SCR_WIDTH/2, SCR_HEIGHT/2)
}

func (arrow *Arrow) Draw(screen *ebiten.Image) {
	screen.DrawImage(arrow.Img, arrow.Op)
}

// Рисуем деления часов - циферблат.
func drawClockPoints(screen *ebiten.Image) {
	var (
		radiusNum      float64 = SCR_WIDTH/2 - 70 // радиус расположения рисок
		radiusPoint    float64
		CenterClockX   float64 = SCR_WIDTH / 2
		CenterClockY   float64 = SCR_HEIGHT / 2
		xPoint, yPoint float64
	)
	bubblesImg := getImage(bubbles)
	smallPointImg := getImage(smallPoint)
	frameSize := 24
	step := 0
	// Риски аналоговых часов
	for a := 0; a < 60; a++ {
		if a%5 == 0 {
			radiusPoint = 12
		} else {
			radiusPoint = 6
		}
		xPoint = CenterClockX + radiusNum*math.Cos(-6.0*float64(a)*(M_PI/180.0)+M_PI/2)
		yPoint = CenterClockY - radiusNum*math.Sin(-6.0*float64(a)*(M_PI/180.0)+M_PI/2)

		if radiusPoint == 12 {
			imgBbl := bubblesImg.SubImage(image.Rect(frameSize*step, 0, frameSize*step+frameSize, 24))
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(xPoint-12, yPoint-12)
			screen.DrawImage(imgBbl.(*ebiten.Image), op)
			step++
		} else {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(xPoint-6, yPoint-6)
			screen.DrawImage(smallPointImg, op)
		}
	}
}

func getImage(arrowByte []byte) *ebiten.Image {
	imgData, _, err := image.Decode(bytes.NewReader(arrowByte))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(imgData)
}
