package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/masfuulaji/go-rpg/entities"
)

type Game struct {
	player       *entities.Player
	TilemapJSON  *TilemapJSON
	TilemapImage *ebiten.Image
	camera       *Camera
	enemies      []*entities.Enemy
	potions      []*entities.Potion
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.X += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.X -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Y -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Y += 2
	}

	for _, sprite := range g.enemies {
		if sprite.FollowsPlayer {
			if sprite.X < g.player.X {
				sprite.X += 1
			}
			if sprite.X > g.player.X {
				sprite.X -= 1
			}
			if sprite.Y < g.player.Y {
				sprite.Y += 1
			}
			if sprite.Y > g.player.Y {
				sprite.Y -= 1
			}
		}
	}

	for _, sprite := range g.potions {
		if g.player.X < sprite.X {
			g.player.Health += sprite.AmtHeal
			fmt.Println("Health:", g.player.Health)
		}
	}

	g.camera.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.camera.Constraint(
		float64(g.TilemapJSON.Layers[0].Width)*16,
		float64(g.TilemapJSON.Layers[0].Height)*16,
		320,
		240,
	)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.TilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(
				g.TilemapImage.SubImage(
					image.Rect(srcX, srcY, srcX+16, srcY+16),
				).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)

	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(
		g.player.Image.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts := ebiten.DrawImageOptions{}

		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)

		screen.DrawImage(
			sprite.Image.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, sprite := range g.potions {
		opts := ebiten.DrawImageOptions{}

		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)

		screen.DrawImage(
			sprite.Image.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImage, _, err := ebitenutil.NewImageFromFile("assets/images/inspector.png")
	if err != nil {
		log.Fatal(err)
	}

	nobleImage, _, err := ebitenutil.NewImageFromFile("assets/images/noble.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImage, _, err := ebitenutil.NewImageFromFile("assets/images/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImage, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(
		&Game{
			player: &entities.Player{
				Sprite: &entities.Sprite{
					Image: playerImage,
					X:     160,
					Y:     220,
				},
				Health: 3,
			},
			enemies: []*entities.Enemy{
				{
					Sprite: &entities.Sprite{
						Image: nobleImage,
						X:     100,
						Y:     100,
					},
					FollowsPlayer: false,
				},
				{
					Sprite: &entities.Sprite{
						Image: nobleImage,
						X:     200,
						Y:     100,
					},
					FollowsPlayer: true,
				},
			},
			potions: []*entities.Potion{
				{
					Sprite: &entities.Sprite{
						Image: potionImage,
						X:     50,
						Y:     200,
					},
					AmtHeal: 2,
				},
			},
			TilemapJSON:  tilemapJSON,
			TilemapImage: tilemapImage,
			camera:       NewCamera(0, 0),
		},
	); err != nil {
		log.Fatal(err)
	}
}
