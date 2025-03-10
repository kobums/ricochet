package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	GRID_SIZE = 33 // 33×33 보드
)

// 33×33 격자에서 벽과 공간을 표현하는 구조체
type Cell struct {
	IsWall bool
}

type Position struct {
	X, Y int
}

type Robot struct {
	Color string
	Pos   Position
}

type GameState struct {
	Board  [GRID_SIZE][GRID_SIZE]Cell // 33x33의 확장된 보드
	Robots map[string]Robot           // 색상별 로봇 위치
	Target Position                   // 목표 위치
}

var game GameState

// 초기 보드 및 로봇 세팅
func initializeGame() {
	// 보드 초기화: 짝수 좌표는 벽, 홀수 좌표는 이동 공간
	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			if x == 0 || x == GRID_SIZE-1 || y == 0 || y == GRID_SIZE-1 {
				game.Board[y][x] = Cell{IsWall: true} // 바깥 벽
			} else if x%2 == 0 || y%2 == 0 {
				game.Board[y][x] = Cell{IsWall: false} // 짝수 좌표는 벽
			} else {
				game.Board[y][x] = Cell{IsWall: false} // 홀수 좌표는 이동 공간
			}
		}
	}

	centerWalls := []struct {
		X, Y int
	}{
		// {2, 12}, {0, 20}, {2, 28}, {4, 6}, {6, 22}, {8, 10}, {10, 6}, {12, 6}, {14, 28},
		{12, 1}, {20, 1}, {28, 3}, {6,5}, {22, 5}, {10, 7}, {6, 9}, {10, 11}, {28, 13}, {22, 15}, {22, 19}, {8, 21}, {12, 23}, {30, 23}, {4, 25}, {18, 25}, {8, 27}, {26, 29}, {8, 31}, {28, 31},
		{29, 2}, {7,4}, {23, 6}, {1, 8}, {5, 8}, {11, 8}, {31, 8}, {9, 12}, {21, 14}, {27, 14},{1, 20}, {23, 20}, {31, 20}, {7, 22}, {11, 22}, {29, 22}, {19, 24}, {5, 26}, {9, 26}, {25, 30},
		{14, 15}, {14, 17}, {18, 15}, {18, 17}, {15, 14}, {17, 14}, {15, 18}, {17, 18}, // centerWall
	}
	for _, w := range centerWalls {
		game.Board[w.Y][w.X] = Cell{IsWall: true}
	}

	// 로봇 초기 위치 설정 (홀수 좌표에서만 배치)
	colors := []string{"R", "B", "G", "Y"}
	game.Robots = make(map[string]Robot)
	rand.Seed(time.Now().UnixNano())

	for _, color := range colors {
		for {
			x, y := rand.Intn(GRID_SIZE/2)*2+1, rand.Intn(GRID_SIZE/2)*2+1 // 홀수 좌표
			if game.Board[y][x].IsWall {
				continue // 벽 위에는 배치 불가능
			}
			game.Robots[color] = Robot{Color: color, Pos: Position{x, y}}
			break
		}
	}

	// 목표 위치 설정 (홀수 좌표에서만 배치)
	for {
		x, y := rand.Intn(GRID_SIZE/2)*2+1, rand.Intn(GRID_SIZE/2)*2+1 // 홀수 좌표
		if game.Board[y][x].IsWall {
			continue
		}
		game.Target = Position{x, y}
		break
	}

	fmt.Println("Game initialized!")
}

// 콘솔에 보드 출력
func printBoard() {
	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			if game.Board[y][x].IsWall {
				fmt.Print("█") // 벽
			} else if y%2 == 0 {
				fmt.Print("-")
			} else if x%2 == 0 {
				fmt.Print("|")
			} else {
				robotFound := false
				for _, robot := range game.Robots {
					if robot.Pos.X == x && robot.Pos.Y == y {
						fmt.Print(string(robot.Color[0])) // 로봇 색상 첫 글자
						robotFound = true
						break
					}
				}
				if !robotFound {
					if game.Target.X == x && game.Target.Y == y {
						fmt.Print("★") // 목표 위치
					} else {
						fmt.Print("·") // 빈 공간
					}
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// 로봇 이동 함수 (33×33 배열 기준)
func moveRobot(robotColor string, direction string) {
	robot, exists := game.Robots[robotColor]
	if !exists {
		fmt.Println("Invalid robot color")
		return
	}

	x, y := robot.Pos.X, robot.Pos.Y
	dx, dy := 0, 0

	switch direction {
	case "up":
		dy = -2
	case "down":
		dy = 2
	case "left":
		dx = -2
	case "right":
		dx = 2
	}

	// 이동 로직: 홀수 좌표에서만 이동하며, 벽 또는 로봇을 만나면 멈춤
	for {
		nx, ny := x+dx, y+dy

		// 보드 밖으로 나가는지 체크
		if nx < 0 || ny < 0 || nx >= GRID_SIZE || ny >= GRID_SIZE {
			break
		}

		// 벽 체크 (짝수 공간에 벽이 있으면 멈춤)
		wx, wy := (x+nx)/2, (y+ny)/2
		if game.Board[wy][wx].IsWall {
			break
		}

		// 다른 로봇과 충돌 체크 (홀수 공간에서 다른 로봇을 만나면 멈춤)
		collision := false
		for _, other := range game.Robots {
			if other.Pos.X == nx && other.Pos.Y == ny {
				collision = true
				break
			}
		}
		if collision {
			break
		}

		// 이동
		x, y = nx, ny
	}

	robot.Pos = Position{x, y}
	game.Robots[robotColor] = robot
}

func main() {
	initializeGame()
	printBoard()

	// 예제: 로봇 이동 (사용자가 직접 입력하도록 변경 가능)
	moveRobot("R", "right")
	printBoard()

	moveRobot("R", "right")
	printBoard()

	moveRobot("G", "up")
	printBoard()

	moveRobot("Y", "down")
	printBoard()

	moveRobot("B", "up")
	printBoard()
}