package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type boundary struct {
	x, y int
}

type block struct {
	id        string
	editorId  int
	tex       *ebiten.Image
	texAnim   bool
	texDivide int
	texOffset *boundary
	hitbox    *boundary
}

type placedBlock struct {
	block block
	pos   *position
	place string
}

func (g *Game) loadBlocks() {
	files, _ := resourcesFS.ReadDir("resources/data/blocks")
	for i, f := range files {
		var result dict
		fileName := fmt.Sprintf("resources/data/blocks/%v", f.Name())
		file, _ := resourcesFS.ReadFile(fileName)
		json.Unmarshal(file, &result)
		fmt.Println(result)

		texPath := fmt.Sprintf("resources/textures/%v.png", result.D("texture").S("src"))
		imageBytes, _ := resourcesFS.ReadFile(texPath)
		decoded, _, _ := image.Decode(bytes.NewReader(imageBytes))
		decodedBlockImage := ebiten.NewImageFromImage(decoded)

		g.blocks = append(g.blocks, &block{editorId: i})
		g.blocks[i].id = result.S("id")
		g.blocks[i].tex = decodedBlockImage
		if result.D("texture").B("animated") {
			g.blocks[i].texAnim = result.D("texture").B("animated")
			g.blocks[i].texDivide = result.D("texture").I("divide")
		}
		g.blocks[i].texOffset = &boundary{x: int(result.D("texture").SD("offset")[0].(float64)), y: int(result.D("texture").SD("offset")[1].(float64))}
		g.blocks[i].hitbox = &boundary{x: int(result.SD("hitbox")[0].(float64)), y: int(result.SD("hitbox")[0].(float64))}
	}
}
