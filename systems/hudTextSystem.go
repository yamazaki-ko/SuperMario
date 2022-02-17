package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	// TextNONE : テキストなし
	TextNONE = 0
	// TextTITLE : タイトル
	TextTITLE = 1
	// TextGOAL  : ゴール
	TextGOAL = 2
	// TextEND : 終了
	TextEND = 3
)

// Text is an entity containing text printed to the screen
type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	textNo   int
	ifMaking bool
}

// HUDTextSystem prints the text to our HUD based on the current state of the game
type HUDTextSystem struct {
	world      *ecs.World
	TextEntity *Text
}

// Update is
func (h *HUDTextSystem) Update(dt float32) {
	if engo.Input.Button("Enter").Down() {
		switch h.TextEntity.textNo {
		case TextTITLE:
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *PlayerSystem:
					sys.playerEntity.ifStart = true
				}
			}
			h.Remove(h.TextEntity.BasicEntity)
			h.TextEntity.textNo = TextNONE
			ifGameOver = false

		case TextGOAL, TextEND:
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *PlayerSystem:
					// リトライ
					sys.PlayerInit(sys.playerEntity)
				}
			}
			h.TextInit(h.TextEntity, TextTITLE)
		}
	}
}

// Remove takes an Entity out of the RenderSystem.
func (h *HUDTextSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// New is
func (h *HUDTextSystem) New(w *ecs.World) {
	h.world = w
	// Entitiy作成
	text := &Text{BasicEntity: ecs.NewBasic()}
	// 初期化
	h.TextInit(text, TextTITLE)
}

// TextInit initializes the value of TextEntity
func (h *HUDTextSystem) TextInit(text *Text, textNo int) {
	// 初期化
	text.textNo = textNo
	text.ifMaking = true
	// SpaceComponent
	TextPositionX := (float32)(0)
	TextPositionY := engo.WindowHeight() - 220
	size := float64(40)
	// RenderComponent
	textDisplay := ""

	// テキストNo毎の処理
	switch textNo {
	case TextTITLE:
		textDisplay = "         GAME START!"
	case TextGOAL:
		textDisplay = "             GOAL!!"
	case TextEND:
		textDisplay = "          GAME OVER"
	}

	// SpaceComponent
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: TextPositionX, Y: TextPositionY},
	}

	// RenderComponent
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.White,
		Size: size,
	}
	fnt.CreatePreloaded()

	text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: textDisplay,
	}

	text.SetShader(common.TextHUDShader)
	text.RenderComponent.SetZIndex(10)

	// Entitiy追加
	h.TextEntity = text
	// SystemにEntity追加
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
		}
	}
}
