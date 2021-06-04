package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type dict map[string]interface{}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampI(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func getNumber() string {
	if inpututil.KeyPressDuration(ebiten.KeyDigit0) == 1 {
		return "0"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit1) == 1 {
		return "1"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit2) == 1 {
		return "2"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit3) == 1 {
		return "3"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit4) == 1 {
		return "4"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit5) == 1 {
		return "5"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit6) == 1 {
		return "6"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit7) == 1 {
		return "7"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit8) == 1 {
		return "8"
	} else if inpututil.KeyPressDuration(ebiten.KeyDigit9) == 1 {
		return "9"
	} else if inpututil.KeyPressDuration(ebiten.KeyBackspace) == 1 {
		return "--"
	}
	return ""
}

func getCursorPositionRelativeToImageWithCameraMovementAndZoom(cursorX, cursorY int, img *ebiten.Image, camera camera) (int, int) {
	imageX, imageY := img.Size()
	relativeCursorX := int((float64(cursorX) - camera.pos.x + float64(imageX)*camera.zoom/2) / camera.zoom)
	relativeCursorY := -int((float64(cursorY) - camera.pos.y - float64(imageY)*camera.zoom/2) / camera.zoom)
	return relativeCursorX, relativeCursorY
}

func (d dict) D(s string) dict {
	return d[s].(map[string]interface{})
}

func (d dict) S(s string) string {
	return d[s].(string)
}

func (d dict) F(s string) float64 {
	return d[s].(float64)
}

func (d dict) I(s string) int {
	return d[s].(int)
}

func (d dict) B(s string) bool {
	return d[s].(bool)
}

func (d dict) SD(s string) []interface{} {
	return d[s].([]interface{})
}
