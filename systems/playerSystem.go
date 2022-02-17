package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

const (
	// MoveDistance : 移動距離
	MoveDistance = 4
	// JumpHeight : ジャンプの高さ
	JumpHeight = 4
	// MaxCount :
	MaxCount = 40
	// PlayerSpriteSheetCell : スプライトシートで使用する最初のセル番号
	PlayerSpriteSheetCell = 8
	// ExtraSizeX : 　プレイヤー画像の余分サイズ
	ExtraSizeX = 8
)

var playerFile = "./Mario/Characters/Mario.png"
var ifGameOver bool

// Player is struct for the PlayerSystem
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// Y初期値
	playerPositionY float32
	// 左右の足の位置
	LeftPositionX  float32
	RightPositionX float32
	// カメラの進んだ距離
	cameraMoveDistance int
	// スプライトシート
	spritesheet *common.Spritesheet
	// 使用中のセル番号
	useCell int
	// ジャンプのカウント数
	jumpCount int
	// 2段ジャンプ用のカウント数
	jumpCount2Step int
	// 頂点までのカウント数
	topCount int
	// 着地点までのカウント数
	bottomCount int
	// ジャンプしているか
	ifJumping bool
	// パイプ上にいるか
	ifOnPipe bool
	// 落下しているか
	ifFalling bool
	// スタートしたか
	ifStart bool
}

// PlayerSystem create a Player to operate
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
}

// Remove removes an Entity from the System
func (ps *PlayerSystem) Remove(ecs.BasicEntity) {
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(ps.playerEntity.BasicEntity)
		}
	}
}

// Update is ran every frame, with `dt` being the time in seconds since the last frame
func (ps *PlayerSystem) Update(dt float32) {
	// 死亡したらリターン
	if ifGameOver {
		return
	}
	// スタートしていなければリターン
	if !ps.playerEntity.ifStart {
		return
	}
	// Goal地点に達したら右移動はしない
	if int(ps.playerEntity.LeftPositionX) >= (TileNum-GoalTileNum+2)*CellWidth16 {
		for _, system := range ps.world.Systems() {
			switch sys := system.(type) {
			case *HUDTextSystem:
				sys.TextInit(sys.TextEntity, TextGOAL)
			}
		}
		ps.Remove(ps.playerEntity.BasicEntity)
	}
	// 落とし穴に落ちる
	if ps.playerEntity.jumpCount == 0 {
		if getMakingInfo(FallPoint, int(ps.playerEntity.LeftPositionX)) && getMakingInfo(FallPoint, int(ps.playerEntity.RightPositionX)) {
			ps.playerEntity.ifFalling = true
			ps.playerEntity.SpaceComponent.Position.Y += MoveDistance
		}
	}
	if ps.playerEntity.SpaceComponent.Position.Y > engo.WindowHeight() {
		ps.PlayerDie()
	}
	if ps.playerEntity.ifFalling {
		return
	}

	// 通常時は動作なし
	if ps.playerEntity.SpaceComponent.Position.Y == ps.playerEntity.playerPositionY {
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell)
	}
	// 土管上にいる時も動作なし
	if ps.playerEntity.ifOnPipe {
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell)
	}

	// プレイヤーを右に移動
	if engo.Input.Button("MoveRight").Down() {
		// 土管位置で土管より下にいる場合
		if getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) && int(ps.playerEntity.SpaceComponent.Position.Y) > int(engo.WindowHeight())-CellHeight16*8 {
			// 右移動できない
		} else {
			// 土管上にいる かつ ジャンプ中でない
			if ps.playerEntity.ifOnPipe && ps.playerEntity.jumpCount == 0 {
				// 土管位置から外れた場合
				if !getMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) && !getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
					ps.playerEntity.ifOnPipe = false
					ps.playerEntity.SpaceComponent.Position.Y = ps.playerEntity.playerPositionY
				}
			}
			// 画面の真ん中より左に位置していれば、カメラを移動せずプレーヤーを移動する
			if int(ps.playerEntity.SpaceComponent.Position.X) < ps.playerEntity.cameraMoveDistance+int(engo.WindowWidth())/2 {
				ps.playerEntity.SpaceComponent.Position.X += MoveDistance
				ps.playerEntity.LeftPositionX += MoveDistance
				ps.playerEntity.RightPositionX += MoveDistance
			} else {
				// 画面の右端に達していなければプレイヤーを移動する
				if int(ps.playerEntity.SpaceComponent.Position.X) < int(engo.WindowWidth())-CellWidth32 {
					ps.playerEntity.SpaceComponent.Position.X += MoveDistance
					ps.playerEntity.LeftPositionX += MoveDistance
					ps.playerEntity.RightPositionX += MoveDistance
				}
				if int(ps.playerEntity.SpaceComponent.Position.X) < TileNum*CellWidth16-int(engo.WindowWidth()/2) {
					// カメラを移動する
					engo.Mailbox.Dispatch(common.CameraMessage{
						Axis:        common.XAxis,
						Value:       MoveDistance,
						Incremental: true,
					})
				}
				ps.playerEntity.cameraMoveDistance += MoveDistance
			}

		}
		// ジャンプ中でない場合
		if ps.playerEntity.jumpCount == 0 {
			switch ps.playerEntity.useCell {
			case 0:
				ps.playerEntity.useCell = 1
			case 1:
				ps.playerEntity.useCell = 2
			case 2:
				ps.playerEntity.useCell = 3
			case 3:
				ps.playerEntity.useCell = 4
			case 4:
				ps.playerEntity.useCell = 0
			}
		} else {
			ps.playerEntity.useCell = 3
		}
		// プレイヤーの動作を変更
		ps.playerEntity.RenderComponent.Drawable = ps.playerEntity.spritesheet.Cell(PlayerSpriteSheetCell + ps.playerEntity.useCell)
	}

	// プレイヤーをジャンプ
	if engo.Input.Button("Jump").JustPressed() {
		// 2段ジャンプ
		if ps.playerEntity.ifJumping {
			if ps.playerEntity.jumpCount < MaxCount/2 {
				ps.playerEntity.jumpCount2Step = ps.playerEntity.jumpCount - 1
			} else {
				ps.playerEntity.jumpCount2Step = MaxCount - (ps.playerEntity.jumpCount - 1)
			}
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = false
		}
		// 初回ジャンプ
		if ps.playerEntity.jumpCount == 0 {
			ps.playerEntity.jumpCount2Step = 0
			ps.playerEntity.jumpCount = 1
			ps.playerEntity.ifJumping = true
		}
		// 土管上からジャンプしていた場合
		if ps.playerEntity.ifOnPipe {
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step + 8
		} else { // 地面からジャンプしていた場合
			ps.playerEntity.bottomCount = 1 + MaxCount + ps.playerEntity.jumpCount2Step
		}
	}

	if ps.playerEntity.jumpCount != 0 {
		ps.playerEntity.jumpCount++
		if ps.playerEntity.jumpCount <= ps.playerEntity.topCount {
			// Up
			ps.playerEntity.SpaceComponent.Position.Y -= JumpHeight
		} else if ps.playerEntity.jumpCount <= ps.playerEntity.bottomCount {
			if ps.playerEntity.SpaceComponent.Position.Y == onPipePositionY {
				// 右足もしくは左足が土管上の場合
				if getMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
					ps.playerEntity.jumpCount = 0
					ps.playerEntity.ifJumping = false
					ps.playerEntity.ifOnPipe = true
				} else {
					// Down
					ps.playerEntity.SpaceComponent.Position.Y += JumpHeight
				}
			} else {
				// Down
				ps.playerEntity.SpaceComponent.Position.Y += JumpHeight
			}

		} else {
			ps.playerEntity.jumpCount = 0
			ps.playerEntity.ifJumping = false
			// 着地点が土管上の場合
			if getMakingInfo(PipePoint, int(ps.playerEntity.LeftPositionX)) || getMakingInfo(PipePoint, int(ps.playerEntity.RightPositionX)) {
				ps.playerEntity.ifOnPipe = true
			} else { // 着地点が地面の場合
				ps.playerEntity.ifOnPipe = false
			}
		}
	}
}

// New is the initialisation of the System
func (ps *PlayerSystem) New(w *ecs.World) {
	//　Worldの追加
	ps.world = w
	//　Entity生成
	player := Player{BasicEntity: ecs.NewBasic()}

	ps.PlayerInit(&player)

	// カメラ設定
	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 3200, Y: 300},
	}
}

// PlayerInit initializes the value of PlayerEntity
func (ps *PlayerSystem) PlayerInit(player *Player) {

	// XY初期値
	PsPositionX := float32(0)
	PsPositionY := engo.WindowHeight() - CellHeight16*6

	// SpaceComponent
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: PsPositionX, Y: PsPositionY},
		Width:    30,
		Height:   30,
	}

	// スプライトシート
	player.spritesheet = common.NewSpritesheetWithBorderFromFile(playerFile, 32, 32, 0, 0)
	// RenderComponent
	player.RenderComponent = common.RenderComponent{
		Drawable: player.spritesheet.Cell(PlayerSpriteSheetCell),
		Scale:    engo.Point{X: 1, Y: 1},
	}
	player.RenderComponent.SetZIndex(5)

	// コンポーネントセット
	ps.playerEntity = player

	// 初期化
	ps.playerEntity.playerPositionY = PsPositionY
	ps.playerEntity.LeftPositionX = PsPositionX + float32(ExtraSizeX)
	ps.playerEntity.RightPositionX = PsPositionX + CellWidth32 - float32(ExtraSizeX)
	ps.playerEntity.ifFalling = false
	ps.playerEntity.ifOnPipe = false
	ps.playerEntity.cameraMoveDistance = 0
	ps.playerEntity.topCount = 1 + MaxCount/2
	ps.playerEntity.bottomCount = 0
	ps.playerEntity.ifStart = false
	ifGameOver = false

	// RenderSystemに追加
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}

	// カメラを移動する
	engo.Mailbox.Dispatch(common.CameraMessage{
		Axis:        common.XAxis,
		Value:       engo.WindowWidth() / 2,
		Incremental: false,
	})
}

// PlayerDie is a function when the Player dies
func (ps *PlayerSystem) PlayerDie() {
	ifGameOver = true
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *HUDTextSystem:
			sys.TextInit(sys.TextEntity, TextEND)
		}
	}
	ps.Remove(ps.playerEntity.BasicEntity)
}
