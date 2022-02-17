package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	// EneymyType0 : パックンフラワー
	EneymyType0 = 0
	// Type0Count :
	Type0Count = 128
	// ExtraSizeXType0 : 余分サイズ
	ExtraSizeXType0 = 6
	// ExtraSizeYType0 : 余分サイズ
	ExtraSizeYType0 = 8
)

var enermyFile = "./Mario/Characters/Enemies.png"

var ifTouched bool

// EnetmyPositionType0 : 敵キャラの位置
var EnetmyPositionType0 []int

// Enermy is struct for the EnermySystem
type Enermy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	count      int
	enermyType int
}

// EnermySystem creates enemies that disturb the player.
type EnermySystem struct {
	world        *ecs.World
	enermyEntity []*Enermy
}

// Remove removes an Entity from the System
func (es *EnermySystem) Remove(ecs.BasicEntity) {
}

// Update is ran every frame, with `dt` being the time in seconds since the last frame
func (es *EnermySystem) Update(dt float32) {

	var playerLeftPositionX float32
	var playerRightPositionX float32
	var playerBottomPositionY float32

	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *PlayerSystem:
			if ifTouched {
				if ifGameOver {
					return
				}
				ifTouched = false
				sys.PlayerDie()
				ifGameOver = true
			}
			playerLeftPositionX = sys.playerEntity.LeftPositionX
			playerRightPositionX = sys.playerEntity.RightPositionX
			playerBottomPositionY = sys.playerEntity.SpaceComponent.Position.Y + float32(CellHeight32)
		}
	}
	if ifGameOver {
		return
	}
	for _, entity := range es.enermyEntity {

		if getMakingInfo(EnetmyPositionType0, int(playerLeftPositionX)) || getMakingInfo(EnetmyPositionType0, int(playerRightPositionX)) {
			if pipePositionY >= playerBottomPositionY && entity.SpaceComponent.Position.Y+ExtraSizeYType0 < playerBottomPositionY {
				ifTouched = true
			}
		}
		if entity.enermyType == EneymyType0 {
			if entity.count < Type0Count {
				entity.SpaceComponent.Position.Y = pipePositionY - float32(entity.count/4)
			} else if entity.count < Type0Count*2 {
				// 一時静止
			} else if entity.count < Type0Count*3 {
				entity.SpaceComponent.Position.Y = pipePositionY - CellHeight32 + float32((entity.count-Type0Count*2)/4)
			} else {
				entity.count = 0
			}
			entity.count++
		}
	}
}

// New is the initialisation of the System
func (es *EnermySystem) New(w *ecs.World) {
	//　Worldの追加
	es.world = w
	ifTouched = false
	ifGameOver = false
	// Enermy配列作成
	Enemies := make([]*Enermy, 0)

	// スプライトシートの作成
	Spritesheet32x32 := common.NewSpritesheetWithBorderFromFile(enermyFile, CellWidth32, CellHeight32, 0, 0)

	//randomNum := rand.Intn(10)
	for i := 0; i <= TileNum; i++ {
		if getMakingInfo(PipePoint, i*CellWidth16) {
			enermy := &Enermy{BasicEntity: ecs.NewBasic()}

			// SpaceComponent
			enermy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(i * CellWidth16), Y: pipePositionY},
			}

			// RenderComponent
			enermy.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet32x32.Cell(7),
				Scale:    engo.Point{X: 1, Y: 1},
			}
			enermy.RenderComponent.SetZIndex(6)

			// 初期化
			enermy.count = 0

			// コンポーネントセット
			Enemies = append(Enemies, enermy)

			// 敵キャラの位置記録
			for j := 0; j < CellWidth32; j++ {
				if j > ExtraSizeXType0 && j < CellWidth32-ExtraSizeXType0 {
					EnetmyPositionType0 = append(EnetmyPositionType0, i*CellWidth16+j)
				}
			}
			i++
		}
	}
	// RenderSystemに追加
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Enemies {
				es.enermyEntity = append(es.enermyEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}
