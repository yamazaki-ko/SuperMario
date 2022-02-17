package systems

import (
	"fmt"
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	// TileNum : Tile数
	TileNum = 200
	// GoalTileNum ： 最終地点からゴールまでのタイル数
	GoalTileNum = 10
	// AroundGoalTileNum ： ゴール付近のタイル数
	AroundGoalTileNum = 30
	// MountTileNum : 山のTile数
	MountTileNum = 5
	// PipeTileNum : 土管のTile数
	PipeTileNum = 1
	// CellWidth16 : 1タイル基準幅(16)
	CellWidth16 = 16
	// CellHeight16 : 1タイル基準高さ(16)
	CellHeight16 = 16
	// CellWidth32 : 1タイル基準幅(32)
	CellWidth32 = 32
	// CellHeight32 : 1タイル基準高さ(32)
	CellHeight32 = 32
	// CellHeight64 : 1タイル基準高さ(64)
	CellHeight64 = 64
	// TileDepth ：深さ
	TileDepth = 4
	// GroundSpriteSheetCell : スプライトシートで使用する地面のセル番号
	GroundSpriteSheetCell = 0
	// CloudSpriteSheetCell : スプライトシートで使用する雲のセル番号
	CloudSpriteSheetCell = 6
	// MountSpriteSheetCell : スプライトシートで使用する山のセル番号
	MountSpriteSheetCell = 11
	// PipeSpriteSheetCell : スプライトシートで使用する土管のセル番号
	PipeSpriteSheetCell = 3
)

var tileFile = "./Mario/Tilesets/OverWorld.png"
var castleFile = "./Mario/Tilesets/Castle.png"

// FallPoint : 落とし穴の位置
var FallPoint []int

// MountPoint : 山の位置
var MountPoint []int

// PipePoint : 土管の位置
var PipePoint []int

// makingxxxx：作成状態（0:作成中でない 1:作成開始 2：それ以外）
var makingFall int
var makingCloud int
var makingMount int
var makingPipe int

// addCell：セル追加値
var addCell int

// Y値
var mountPositionY float32
var pipePositionY float32
var onPipePositionY float32
var castleositionY float32

// Tile is Eintity for the TileSystem
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// TileSystem builds a game background
type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
}

// Remove removes an Entity from the System
func (*TileSystem) Remove(ecs.BasicEntity) {}

// Update is ran every frame, with `dt` being the time in seconds since the last frame
func (ts *TileSystem) Update(dt float32) {}

// New is the initialisation of the System
func (ts *TileSystem) New(w *ecs.World) {
	//　Worldの追加
	ts.world = w
	// スプライトシートの作成
	Spritesheet16x16 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth16, CellHeight16, 0, 0)
	Spritesheet32x32 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth32, CellHeight32, 0, 0)
	Spritesheet16x64 := common.NewSpritesheetWithBorderFromFile(tileFile, CellWidth16, CellHeight64, 0, 0)

	// 初期化
	makingFall = 0
	addCell = 0
	cloudHeight := 0
	mountPositionY = engo.WindowHeight() - CellHeight16*7
	pipePositionY = engo.WindowHeight() - CellHeight16*6
	onPipePositionY = engo.WindowHeight() - CellHeight16*8
	castleositionY = engo.WindowHeight() - CellHeight16*9

	// Tile配列作成
	Tiles := make([]*Tile, 0)

	for i := 0; i <= TileNum; i++ {
		// ----------------------- //
		// ------- 地面の作成 ------ //
		// ----------------------- //
		// Start付近とGoal付近に落とし穴は作らない
		if i >= 10 && i < TileNum-AroundGoalTileNum {
			randomNum := rand.Intn(10)
			if randomNum == 0 {
				makingFall = 1
			} else {
				// 最低タイル2つ分は落とし穴を作成する
				if makingFall == 1 {
					makingFall = 2
				} else {
					makingFall = 0
				}
			}
		}
		if makingFall != 0 {
			for j := 0; j < CellWidth16; j++ {
				// 落とし穴の位置記録
				FallPoint = append(FallPoint, i*CellWidth16+j)
			}
		} else {
			for j := 0; j < TileDepth; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				// SpaceComponent
				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth16), Y: float32(int(engo.WindowHeight()) - (j+1)*CellHeight16)},
				}
				// RenderComponent
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet16x16.Cell(GroundSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				// コンポーネントセット
				Tiles = append(Tiles, tile)
			}

		}
		// ----------------------- //
		// ------- 雲の作成 ------- //
		// ----------------------- //
		if makingCloud == 0 {
			randomNum := rand.Intn(12)
			if randomNum < 3 {
				makingCloud = 1
				cloudHeight = randomNum
			}
		}
		if makingCloud != 0 {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			j := float32(0)
			// 2つ目の雲作成中の場合
			if makingCloud > 2 {
				j = float32(i) - 0.5
			} else {
				j = float32(i)
			}

			// SpaceComponent
			tile.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(j * CellWidth32), Y: float32(int(engo.WindowHeight()/3) - cloudHeight*CellHeight16)},
			}

			// RenderComponent
			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet32x32.Cell(CloudSpriteSheetCell + addCell),
				Scale:    engo.Point{X: 1, Y: 1},
			}
			tile.RenderComponent.SetZIndex(float32(makingCloud))

			// コンポーネントセット
			Tiles = append(Tiles, tile)

			switch makingCloud {
			case 1:
				makingCloud++
				addCell = 1
				break
			case 2:
				makingCloud++
				addCell = 1
				break
			default:
				makingCloud = 0
				addCell = 0
				break
			}
		}
	}
	for i := 0; i <= TileNum; i++ {
		// ----------------------- //
		// ------- 山の作成 ------- //
		// ----------------------- //
		makingMount = 1
		for j := 0; j < MountTileNum+2; j++ {
			// 山を作成できる十分なスペースがない場合
			if getMakingInfo(FallPoint, (i+j)*CellWidth16) {
				makingMount = 0
			}
		}
		if makingMount != 0 && i < TileNum-AroundGoalTileNum {
			for j := 0; j < MountTileNum; j++ {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				// SpaceComponent
				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32((i + j) * CellWidth16), Y: mountPositionY},
				}

				// RenderComponent
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet16x64.Cell(MountSpriteSheetCell + j),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(0)

				// コンポーネントセット
				Tiles = append(Tiles, tile)

				// 山の位置記録
				MountPoint = append(MountPoint, (i+j)*CellWidth16)
			}
			// ランダムな値をインクリメント
			i = i + 20
		}
	}
	for i := 0; i <= TileNum; i++ {
		// ------------------------ //
		// ------- 土管の作成 ------- //
		// ------------------------ //
		makingPipe = 1
		for j := 0; j < PipeTileNum+2; j++ {
			// 土管を作成できる十分なスペースがない、もしくは山を作成している場合
			if getMakingInfo(FallPoint, (i+j)*CellWidth16) || getMakingInfo(MountPoint, (i+j)*CellWidth16) {
				makingPipe = 0
			}
		}
		// Start付近とGoal付近は土管は作らない
		if i >= 10 && i < TileNum-AroundGoalTileNum {
			if makingPipe != 0 {
				tile := &Tile{BasicEntity: ecs.NewBasic()}

				// SpaceComponent
				tile.SpaceComponent = common.SpaceComponent{
					Position: engo.Point{X: float32(i * CellWidth16), Y: pipePositionY},
				}

				// RenderComponent
				tile.RenderComponent = common.RenderComponent{
					Drawable: Spritesheet32x32.Cell(PipeSpriteSheetCell),
					Scale:    engo.Point{X: 1, Y: 1},
				}
				tile.RenderComponent.SetZIndex(7)

				// コンポーネントセット
				Tiles = append(Tiles, tile)

				// 土管の位置記録
				for j := 0; j < CellWidth32; j++ {
					PipePoint = append(PipePoint, i*CellWidth16+j)
				}
				// ランダムな値をインクリメント
				i = i + 30
			}
		}
	}
	// ----------------------- //
	// ------- 城の作成 ------- //
	// ----------------------- //
	tile := &Tile{BasicEntity: ecs.NewBasic()}

	// SpaceComponent
	tile.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32((TileNum - GoalTileNum) * CellWidth16), Y: castleositionY},
	}

	// 画像の読み込み
	texture, err := common.LoadedSprite(castleFile)
	if err != nil {
		fmt.Println("Unable to load texture: " + castleFile + "：" + err.Error())
	}
	// RenderComponent
	tile.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	tile.RenderComponent.SetZIndex(0)

	// コンポーネントセット
	Tiles = append(Tiles, tile)

	// RenderSystemに追加
	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}

// getMakingInfo ： 対象位置に含まれているか
func getMakingInfo(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
