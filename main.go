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

var oneBoard = []Position{
	{12, 1}, {6, 5}, {10, 7}, {6, 9}, {10, 11},
	{7, 4}, {1, 8}, {5, 8}, {11, 8}, {9, 12},
}

var twoBoard = []Position{
	{1, 12}, {3, 4}, {5, 10}, {13, 4}, {15, 10},
	{2, 3}, {6, 9}, {8, 1}, {14, 5}, {14, 11},
}

var threeBoard = []Position{
	{12, 1}, {6, 5}, {10, 7}, {6, 9}, {10, 11},
	{7, 4}, {1, 8}, {5, 8}, {11, 8}, {9, 12},
}

var fourBoard = []Position{
	{4, 1}, {6, 3}, {2, 9}, {10,  13}, {14, 7},
	{1, 12}, {3, 10}, {7, 2}, {9, 12}, {13, 8},
}

var centerWall = []Position{
	{14, 15}, {14, 17}, {18, 15}, {18, 17}, {15, 14}, {17, 14}, {15, 18}, {17, 18},
}

var centerBoard = []Position{
	{15, 15}, {15, 17}, {17, 15}, {17, 17},
}

func transformPosition(p Position, index int) Position {
	switch index {
	case 1:
		return p // 1사분면: 그대로 유지
	case 2:
		return Position{GRID_SIZE - p.Y - 1, p.X} // 2사분면 변환
	case 3:
		return Position{p.Y, GRID_SIZE - p.X - 1} // 3사분면 변환
	case 4:
		return Position{GRID_SIZE - p.X - 1, GRID_SIZE - p.Y - 1} // 4사분면 변환
	default:
		return p // 기본값: 변화 없음
	}
}

// 보드를 4개의 사분면으로 나눈 후 섞기 (벽도 함께 이동)
func shuffleQuadrants() {
	rand.Seed(time.Now().UnixNano())

	// 사분면의 순서를 랜덤하게 섞음
	quadrantOrder := []int{1, 2, 3, 4}
	rand.Shuffle(len(quadrantOrder), func(i, j int) { quadrantOrder[i], quadrantOrder[j] = quadrantOrder[j], quadrantOrder[i] })

	newWalls := make([]Position, len(oneBoard)*4)

	// 기존 사분면 벽을 변환하여 newWalls에 추가
	for i, newQ := range quadrantOrder {
		var boardToUse []Position

		// 각 사분면에 해당하는 벽 데이터 선택
		switch newQ {
		case 1:
			boardToUse = oneBoard
		case 2:
			boardToUse = twoBoard
		case 3:
			boardToUse = threeBoard
		case 4:
			boardToUse = fourBoard
		}

		// 변환된 벽 좌표 추가
		for _, p := range boardToUse {
			newWalls = append(newWalls, transformPosition(p, i+1))
		}
	}

	// 기존 보드를 초기화하고 새로운 벽 설정
	for y := 1; y < GRID_SIZE-1; y++ {
		for x := 1; x < GRID_SIZE-1; x++ {
			game.Board[y][x] = Cell{IsWall: false}
		}
	}
	for _, w := range newWalls {
		game.Board[w.Y][w.X] = Cell{IsWall: true}
	}
}


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

	// 기본 벽 추가 후 사분면 섞기
	shuffleQuadrants()

	for _, w := range centerWall {
		game.Board[w.Y][w.X] = Cell{IsWall: true}
	}

	for _, w := range centerBoard {
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