package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	pb "github.com/KatyaPark11/Sudoku-Golang/generated/sudoku"
	"google.golang.org/grpc"
)

type sudokuServer struct {
	pb.UnimplementedSudokuServiceServer
}

const N = 9

// isSafe проверяет, безопасно ли поставить число num в позицию (row, col) на доске
func isSafe(board [N][N]int, row, col, num int) bool {
	for x := range N {
		if board[row][x] == num || board[x][col] == num ||
			board[(row/3)*3+x/3][(col/3)*3+x%3] == num {
			return false
		}
	}
	return true
}

// hiddenSingles ищет скрытые одиночки и заполняет их
func hiddenSingles(board [N][N]int) bool {
	found := false

	for num := 1; num <= 9; num++ {
		// Проверка боксов
		for boxRow := range 3 {
			for boxCol := range 3 {
				startRow := boxRow * 3
				startCol := boxCol * 3

				var possiblePositions []struct{ row, col int }

				for i := range 3 {
					for j := range 3 {
						row := startRow + i
						col := startCol + j
						if board[row][col] == 0 && isSafe(board, row, col, num) {
							possiblePositions = append(possiblePositions, struct{ row, col int }{row, col})
						}
					}
				}

				if len(possiblePositions) == 1 {
					r := possiblePositions[0].row
					c := possiblePositions[0].col
					board[r][c] = num
					found = true
				}
			}
		}

		// Проверка строк и столбцов
		for i := range N {
			var possibleRowPos = -1
			var possibleColPos = -1

			for j := range N {
				if board[i][j] == 0 && isSafe(board, i, j, num) {
					if possibleRowPos == -1 {
						possibleRowPos = j
					} else {
						possibleRowPos = -2 // больше одной позиции
						break
					}
				}
			}
			if possibleRowPos >= 0 {
				board[i][possibleRowPos] = num
				found = true
			}

			for j := range N {
				if board[j][i] == 0 && isSafe(board, j, i, num) {
					if possibleColPos == -1 {
						possibleColPos = j
					} else {
						possibleColPos = -2 // больше одной позиции
						break
					}
				}
			}
			if possibleColPos >= 0 {
				board[possibleColPos][i] = num
				found = true
			}
		}
	}

	return found
}

// backtrackSolve решает судоку методом проб и ошибок (backtracking)
func backtrackSolve(board *[N][N]int, steps *[][N][N]int, isSteps bool) bool {
	for row := range N {
		for col := range N {
			if board[row][col] == 0 {
				for num := 1; num <= 9; num++ {
					if isSafe(*board, row, col, num) {
						board[row][col] = num
						if isSteps {
							saveStep(*board, steps)
						}
						if backtrackSolve(board, steps, isSteps) {
							return true
						}
						board[row][col] = 0 // откат (backtracking)
					}
				}
				return false // не удалось поставить ни одного числа — назад по цепочке
			}
		}

	}
	return true // все клетки заполнены успешно
}

// saveStep сохраняет текущий шаг решения в список steps
func saveStep(currentBoard [N][N]int, steps *[][N][N]int) {
	stepCopy := currentBoard // копируем текущую доску (по значению)
	*steps = append(*steps, stepCopy)
}

// parseBoardFromString преобразует строку в доску [N][N]int
func parseBoardFromString(s string) ([N][N]int, error) {
	var board [N][N]int
	if len(s) != N*N {
		return board, fmt.Errorf("длина строки должна быть %d символов", N*N)
	}
	for i, ch := range s {
		row := i / N
		col := i % N
		if ch < '0' || ch > '9' {
			return board, fmt.Errorf("недопустимый символ: %c", ch)
		}
		board[row][col] = int(ch - '0')
	}
	return board, nil
}

// solveSudoku принимает начальную доску и флаг необходимости сохранять шаги,
// возвращает строку с решением и шагами.
func solveSudoku(initialBoard string, isSteps bool) (string, string, error) {

	board, err := parseBoardFromString(initialBoard)
	if err != nil {
		return "", "", err
	}

	var steps [][N][N]int

	strategyApplied := true

	for strategyApplied {
		strategyApplied = false

		if hiddenSingles(board) {
			strategyApplied = true
			if isSteps {
				saveStep(board, &steps)
			}
		}

		// Можно добавить другие стратегии по мере необходимости...
	}

	// Продолжаем решение методом проб и ошибок (backtracking)
	success := backtrackSolve(&board, &steps, isSteps)
	if !success {
		return "", "", fmt.Errorf("не удалось решить судоку")
	}

	// Преобразуем финальную доску и шаги в строки для вывода или передачи.
	solutionStr := boardToString(board)
	stepsStr := stepsToString(steps)

	return solutionStr, stepsStr, nil
}

// boardToString преобразует доску в строку, выводя только цифры подряд.
func boardToString(board [N][N]int) string {
	var sb strings.Builder
	for _, row := range board {
		for _, val := range row {
			sb.WriteString(fmt.Sprint(val))
		}
	}
	return sb.String()
}

// stepsToString преобразует список шагов в строку, выводя только цифры подряд.
func stepsToString(steps [][N][N]int) string {
	var sb strings.Builder
	for _, step := range steps {
		for _, row := range step {
			for _, val := range row {
				sb.WriteString(fmt.Sprint(val))
			}
		}
	}
	return sb.String()
}

func (s *sudokuServer) Solve(ctx context.Context, req *pb.SudokuRequest) (*pb.SudokuResponse, error) {

	initialBoard := req.Puzzle
	isSteps := req.IsSteps

	solutionStr, stepsStr, err := solveSudoku(initialBoard, isSteps)

	if err != nil {
		return &pb.SudokuResponse{Solution: ""}, err
	} else if isSteps {
		return &pb.SudokuResponse{Solution: stepsStr}, nil
	} else {
		return &pb.SudokuResponse{Solution: solutionStr}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSudokuServiceServer(grpcServer, &sudokuServer{})

	log.Println("Sudoku service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
