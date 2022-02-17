package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/yamazaki-ko/SuperMario/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type myScene struct{}

func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	engo.Files.Load("./Mario/Characters/Mario.png")
	engo.Files.Load("./Mario/Characters/Enemies.png")
	engo.Files.Load("./assets/Mario/Misc/Items.png")
	engo.Files.Load("./Mario/Tilesets/OverWorld.png")
	engo.Files.Load("./Mario/Tilesets/Castle.png")
	common.SetBackground(color.RGBA{120, 226, 250, 3})
}

// Setup is called before the main loop starts.
// It allows you to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	// キーボード設定
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("Jump", engo.KeySpace)
	engo.Input.RegisterButton("Enter", engo.KeyEnter)
	// フォント設定
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))

	// World設定
	world, _ := u.(*ecs.World)

	// Systemの追加
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.TileSystem{})
	world.AddSystem(&systems.PlayerSystem{})
	world.AddSystem(&systems.EnermySystem{})
	world.AddSystem(&systems.HUDTextSystem{})
}

func main() {
	fmt.Printf("hello, world\n")
	opts := engo.RunOptions{
		Title:          "SuperMario",
		Width:          480,
		Height:         320,
		StandardInputs: true,
		NotResizable:   true,
	}
	fmt.Println("SuperMario Start")
	engo.Run(opts, &myScene{})
}

func (*myScene) Exit() {
	engo.Exit()
}
