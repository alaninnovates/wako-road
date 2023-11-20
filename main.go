package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EnemyCar struct {
	texture rl.Texture2D
	y       float32
	lane    int32
}

func main() {
	rl.InitWindow(450, 800, "raylib")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	carTextures := make([]rl.Texture2D, 0)
	for i := 1; i <= 4; i++ {
		carTexture := rl.LoadTexture(fmt.Sprintf("resources/car%d.png", i))
		carTextures = append(carTextures, carTexture)
		defer rl.UnloadTexture(carTexture)
	}
	coinTexture := rl.LoadTexture("resources/coin.png")

	carTexture := carTextures[0]
	enemyCars := make([]EnemyCar, 0)
	coins := make([]rl.Vector2, 0)

	carX := float32(rl.GetScreenWidth()/2) - float32(carTexture.Width)*0.1/2
	carY := float32(rl.GetScreenHeight()) - float32(carTexture.Height)*0.1 - 10
	carLane := 2
	carSpawnDelay := 0
	coinSpawnDelay := 0
	score := 0
	coinsCount := 0

	// time elapsed to determine speed of cars
	timeElapsed := 0

	gameStatus := "playing"
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if gameStatus == "playing" {
			if rl.IsKeyPressed(rl.KeyRight) && carLane < 4 {
				carX += float32(rl.GetScreenWidth() / 5)
				carLane += 1
			}
			if rl.IsKeyPressed(rl.KeyLeft) && carLane > 0 {
				carX -= float32(rl.GetScreenWidth() / 5)
				carLane -= 1
			}
			drawTracks()
			rl.DrawTextureEx(carTexture, rl.Vector2{X: carX, Y: carY}, 0, 0.1, rl.White)

			// score
			rl.DrawText(
				fmt.Sprintf("Score: %d", score),
				int32(rl.GetScreenWidth()/2)-rl.MeasureText(fmt.Sprintf("Score: %d", score), 20)/2,
				20,
				20,
				rl.Yellow,
			)
			rl.DrawText(
				fmt.Sprintf("Coins: %d", coinsCount),
				10,
				10,
				20,
				rl.Black,
			)
			score += 1
			timeElapsed += 1

			// enemy cars
			carSpawnDelay += 1
			if carSpawnDelay > 40-timeElapsed/100 {
				carSpawnDelay = 0
				enemyCar := EnemyCar{
					texture: carTextures[rl.GetRandomValue(0, 3)],
					y:       0,
					lane:    rl.GetRandomValue(0, 4),
				}
				enemyCars = append(enemyCars, enemyCar)
			}
			// coins
			coinSpawnDelay += 1
			if coinSpawnDelay > 100-timeElapsed/100 {
				coinSpawnDelay = 0
				coin := rl.Vector2{
					X: float32(
						rl.GetScreenWidth()/2,
					) - float32(
						coinTexture.Width,
					)*0.1/2 - float32(
						rl.GetScreenWidth()/5,
					)*(float32(rl.GetRandomValue(0, 4))-2),
					Y: 0,
				}
				coins = append(coins, coin)
			}
			for i := 0; i < len(enemyCars); i++ {
				enemyCars[i].y += float32(5 + timeElapsed/100)
				x := float32(
					rl.GetScreenWidth()/2,
				) - float32(
					carTexture.Width,
				)*0.1/2 - float32(
					rl.GetScreenWidth()/5,
				)*(float32(enemyCars[i].lane)-2)
				rl.DrawTextureEx(enemyCars[i].texture, rl.Vector2{X: x, Y: enemyCars[i].y}, 0, 0.1, rl.White)

				// collisions
				carRect := rl.Rectangle{
					X:      carX,
					Y:      carY,
					Width:  float32(carTexture.Width) * 0.1,
					Height: float32(carTexture.Height) * 0.1,
				}
				enemyCarRect := rl.Rectangle{
					X:      x,
					Y:      enemyCars[i].y,
					Width:  float32(carTexture.Width) * 0.1,
					Height: float32(carTexture.Height) * 0.1,
				}
				if rl.CheckCollisionRecs(carRect, enemyCarRect) {
					gameStatus = "gameover"
				}
			}
			for i := 0; i < len(coins); i++ {
				coins[i].Y += float32(5 + timeElapsed/100)
				rl.DrawTextureEx(coinTexture, coins[i], 0, 0.08, rl.White)
				if rl.CheckCollisionRecs(
					rl.Rectangle{
						X:      carX,
						Y:      carY,
						Width:  float32(carTexture.Width) * 0.1,
						Height: float32(carTexture.Height) * 0.1,
					},
					rl.Rectangle{
						X:      coins[i].X,
						Y:      coins[i].Y,
						Width:  float32(coinTexture.Width) * 0.08,
						Height: float32(coinTexture.Height) * 0.08,
					},
				) {
					coins = append(coins[:i], coins[i+1:]...)
					coinsCount += 1
				}
			}
		} else if gameStatus == "gameover" {
			if rl.IsKeyPressed(rl.KeyR) {
				gameStatus = "playing"
				enemyCars = make([]EnemyCar, 0)
				coins = make([]rl.Vector2, 0)
				carX = float32(rl.GetScreenWidth()/2) - float32(carTexture.Width)*0.1/2
				carY = float32(rl.GetScreenHeight()) - float32(carTexture.Height)*0.1 - 10
				carLane = 2
				carSpawnDelay = 0
				coinSpawnDelay = 0
				score = 0
				timeElapsed = 0
			}
			rl.DrawText(
				"Game Over",
				int32(rl.GetScreenWidth()/2)-rl.MeasureText("Game Over", 40)/2,
				int32(rl.GetScreenHeight()/2),
				40,
				rl.Black,
			)
			rl.DrawText(
				fmt.Sprintf("Score: %d", score),
				int32(rl.GetScreenWidth()/2)-rl.MeasureText(fmt.Sprintf("Score: %d", score), 20)/2,
				int32(rl.GetScreenHeight()/2)+40,
				20,
				rl.Black,
			)
			rl.DrawText(
				"Press R to restart",
				int32(rl.GetScreenWidth()/2)-rl.MeasureText("Press R to restart", 20)/2,
				int32(rl.GetScreenHeight()/2)+80,
				20,
				rl.Black,
			)
		}

		rl.EndDrawing()
	}
}

func drawTracks() {
	// background: grey
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.LightGray)
	// borders
	rl.DrawRectangle(0, 0, 10, int32(rl.GetScreenHeight()), rl.Gray)
	rl.DrawRectangle(int32(rl.GetScreenWidth())-10, 0, 10, int32(rl.GetScreenHeight()), rl.Gray)
	// lanes
	drawDottedTrackLine(int32(rl.GetScreenWidth()/5) - 2)
	drawDottedTrackLine(int32(rl.GetScreenWidth()/5)*2 - 2)
	drawDottedTrackLine(int32(rl.GetScreenWidth()/5)*3 - 2)
	drawDottedTrackLine(int32(rl.GetScreenWidth()/5)*4 - 2)
}

func drawDottedTrackLine(x int32) {
	// draw dotted line that is height of screen
	for i := 0; i < int(rl.GetScreenHeight()); i += 10 {
		rl.DrawRectangle(x, int32(i), 5, 5, rl.Gray)
	}
}
