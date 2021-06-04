package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"strconv"
)

const (
	sWidth  = 1280
	sHeight = 640
)

//go:embed resources
var resourcesFS embed.FS

var (
	normalText         font.Face
	menuImage          = ebiten.NewImage(sWidth, sHeight)
	decodedMenuImage   *ebiten.Image
	blocksImage        *ebiten.Image
	gridImage          *ebiten.Image
	gridImageProcessed = ebiten.NewImage(sWidth, sHeight)
	alphaBlockImage    = ebiten.NewImage(16, 16)
	blocksMenu         = ebiten.NewImage(40, sHeight)
	fullImage          = ebiten.NewImage(sWidth, sHeight)
	messageSizeMenu    [2]string
	selectedSizeMenu   int
)

type position struct {
	x, y float64
}

type camera struct {
	pos  *position
	zoom float64
}

type Mode int

const (
	Init Mode = iota
	Menu
	InitiateGrid
	Play
)

type Game struct {
	mode          Mode
	camera        camera
	selectedBlock int
	blocks        []*block
	world         world
}

type world struct {
	name    string
	levelId string
	maps    []string
}

func (g *Game) Update() error {
	if g.mode == Init {
		g.loadBlocks()
		imageBytes, _ := resourcesFS.ReadFile("resources/textures/gui/menu.png")
		decoded, _, _ := image.Decode(bytes.NewReader(imageBytes))
		decodedMenuImage = ebiten.NewImageFromImage(decoded)
		momImgX, momImgY := decodedMenuImage.Size()
		opM := &ebiten.DrawImageOptions{}
		opM.GeoM.Translate(float64((sWidth-momImgX)/2), float64((sHeight-momImgY)/2))
		menuImage.DrawImage(decodedMenuImage, opM)
		messageSizeMenu[0] = ""
		messageSizeMenu[1] = ""
		g.camera.zoom = 1.0
		g.camera.pos = &position{x: 0, y: 0}
		g.mode++
	} else if g.mode == Menu {
		cx, cy := ebiten.CursorPosition()
		if cx > 500 && cy > 250 && cx <= 860 && cy <= 300 && inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 1 {
			selectedSizeMenu = 0
		} else if cx > 500 && cy > 325 && cx <= 860 && cy <= 440 && inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 1 {
			selectedSizeMenu = 1
		}
		num := getNumber()
		if num == "--" {
			if len(messageSizeMenu[selectedSizeMenu]) >= 1 {
				messageSizeMenu[selectedSizeMenu] = messageSizeMenu[selectedSizeMenu][:len(messageSizeMenu[selectedSizeMenu])-1]
			}
		} else if len(num) > 0 {
			messageSizeMenu[selectedSizeMenu] += num
		}

		momImgX, momImgY := decodedMenuImage.Size()
		opM := &ebiten.DrawImageOptions{}
		opM.GeoM.Translate(float64((sWidth-momImgX)/2), float64((sHeight-momImgY)/2))
		menuImage.DrawImage(decodedMenuImage, opM)

		mes0, _ := strconv.Atoi(messageSizeMenu[0])
		mes1, _ := strconv.Atoi(messageSizeMenu[1])
		mes0 = clampI(mes0*16, 0, 16384)
		mes1 = clampI(mes1*16, 0, 8192)

		text.Draw(menuImage, strconv.Itoa(mes0/16), normalText, 512, 292, color.Black)
		text.Draw(menuImage, strconv.Itoa(mes1/16), normalText, 512, 408, color.Black)

		text.Draw(menuImage, fmt.Sprintf("blocks = %vpx", mes0), normalText, 870, 292, color.Black)
		text.Draw(menuImage, fmt.Sprintf("blocks = %vpx", mes1), normalText, 870, 408, color.Black)

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gridX, _ := strconv.Atoi(messageSizeMenu[0])
			gridY, _ := strconv.Atoi(messageSizeMenu[1])
			if gridX > 0 && gridY > 0 {
				gridImage = ebiten.NewImage(clampI(gridX*16, 16, 16384), clampI(gridY*16, 16, 8192))
				g.mode++
			}
		}
	} else if g.mode == InitiateGrid {
		imageBytes, _ := resourcesFS.ReadFile("resources/textures/editor/grid.png")
		decoded, _, _ := image.Decode(bytes.NewReader(imageBytes))
		decodedGrid := ebiten.NewImageFromImage(decoded)
		opM := &ebiten.DrawImageOptions{}
		gridX, gridY := gridImage.Size()
		blocksImage = ebiten.NewImage(gridX, gridY)
		gridDecodedX, gridDecodedY := decodedGrid.Size()
		for i := 0; i <= gridY/16; i++ {
			for j := 0; j <= gridX/16; j++ {
				moveX := gridDecodedX
				if j == 0 {
					moveX = 0
				}
				opM.GeoM.Translate(float64(moveX), 0)
				gridImage.DrawImage(decodedGrid, opM)
			}
			opM.GeoM.Translate(float64(-gridX), float64(gridDecodedY))
		}
		g.camera.pos = &position{x: float64(sWidth / 2), y: float64(sHeight / 2)}
		g.mode++
	} else if g.mode == Play {
		if inpututil.KeyPressDuration(ebiten.KeyW) > 1 {
			g.camera.pos.y += 2
		}
		if inpututil.KeyPressDuration(ebiten.KeyS) > 1 {
			g.camera.pos.y -= 2
		}
		if inpututil.KeyPressDuration(ebiten.KeyA) > 1 {
			g.camera.pos.x += 2
		}
		if inpututil.KeyPressDuration(ebiten.KeyD) > 1 {
			g.camera.pos.x -= 2
		}
		_, mouseScrolly := ebiten.Wheel()
		g.camera.zoom = clamp(g.camera.zoom+mouseScrolly/10, 0.1, 10.0)

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.selectedBlock = clampI(g.selectedBlock-1, 0, 10)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.selectedBlock = clampI(g.selectedBlock+1, 0, 10)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()
	if g.mode == Menu {
		screen.DrawImage(menuImage, nil)
	} else if g.mode == Play {
		fullImage.Clear()

		gridImageProcessed.Clear()
		optGrid := &ebiten.DrawImageOptions{}
		gridImageX, gridImageY := gridImage.Size()
		optGrid.GeoM.Translate(-float64(gridImageX)/2, -float64(gridImageY)/2)
		optGrid.GeoM.Scale(g.camera.zoom, g.camera.zoom)
		optGrid.GeoM.Translate(g.camera.pos.x, g.camera.pos.y)
		gridImageProcessed.DrawImage(gridImage, optGrid)
		fullImage.DrawImage(gridImageProcessed, nil)

		relCX, relCY := getCursorPositionRelativeToImageWithCameraMovementAndZoom(cx, cy, gridImage, g.camera)
		cIx := clampI(relCX, 0, gridImageX-1)
		cIy := clampI(relCY, 0, gridImageY-1)

		alphaBlockImage.Clear()
		optAlphaBlocks := &ebiten.DrawImageOptions{}
		alphaBlockImageX, alphaBlockImageY := alphaBlockImage.Size()
		optAlphaBlocks.GeoM.Translate(-float64(alphaBlockImageX)/2, -float64(alphaBlockImageY)/2)
		optAlphaBlocks.GeoM.Scale(g.camera.zoom, g.camera.zoom)
		tranX := g.camera.pos.x + (float64(alphaBlockImageX/2)-float64(gridImageX/2)+float64(cIx/16)*16)*g.camera.zoom
		tranY := g.camera.pos.y + (-float64(alphaBlockImageY/2)+float64(gridImageY/2)-float64(cIy/16)*16)*g.camera.zoom
		optAlphaBlocks.GeoM.Translate(tranX, tranY)
		alphaBlockImage.Fill(color.RGBA{R: 122, G: 122, B: 122, A: 122})
		fullImage.DrawImage(alphaBlockImage, optAlphaBlocks)

		fullImage.DrawImage(blocksImage, nil)

		optBlocksMenu := &ebiten.DrawImageOptions{}
		blocksMenu.Fill(color.RGBA{R: 122, G: 122, B: 122, A: 122})
		optBlocksMenu.GeoM.Translate(sWidth-40, 0)
		fullImage.DrawImage(blocksMenu, optBlocksMenu)

		text.Draw(fullImage, fmt.Sprintf("(%v, %v)", cIx, cIy), normalText, cx, cy+15, color.RGBA{R: 255, G: 150, A: 255})
		text.Draw(fullImage, fmt.Sprintf("(%v, %v)", cIx/16, cIy/16), normalText, cx, cy+30, color.RGBA{R: 255, G: 150, A: 255})
		text.Draw(fullImage, fmt.Sprintf("(%v, %v)", g.camera.pos.x, g.camera.pos.y), normalText, cx, cy+45, color.RGBA{R: 255, G: 150, A: 255})

		screen.DrawImage(fullImage, nil)
	}

	text.Draw(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()), normalText, 0, sHeight-30, color.White)
	text.Draw(screen, fmt.Sprintf("(%v, %v)", cx, cy), normalText, cx, cy, color.RGBA{R: 255, G: 150, A: 255})
}

func (g *Game) Layout(int, int) (int, int) {
	return sWidth, sHeight
}

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 50
	normalText, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(sWidth, sHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("PlatformatorMapEditor")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
