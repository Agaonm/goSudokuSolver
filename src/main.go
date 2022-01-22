package main

import (
	"fmt"

	//"os"
	"strconv"
	"syscall/js"
)

var grid [9][9]int
var solvedArray [81]string
var document js.Value

func clearBoard(this js.Value, i []js.Value) interface{} {
	for i := 0; i < 81; i++ {
		document.Call("getElementById", i).Set("value", "")
	}
	return nil
}

func checkSudoku(this js.Value, i []js.Value) interface{} {
	var value1 [81]string
	for i := 0; i < 81; i++ {
		s := strconv.Itoa(i)
		value1[i] = document.Call("getElementById", s).Get("value").String()
	}
	splitArray(value1)
	if backtrack(&grid) {
		fmt.Println("Solved!")
		createArray(&grid)

		for i := 0; i < 81; i++ {
			document.Call("getElementById", i).Set("value", solvedArray[i])
		}

	} else {
		fmt.Println("INVALID CANT BE SOLVED")
	}
	return nil
}

func splitArray(sudoku1d [81]string) {
	i := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			s, err := strconv.Atoi(sudoku1d[i])
			if err != nil {
				s = 0
			}
			grid[row][col] = s
			i++
		}
	}
}

func createArray(grid *[9][9]int) {
	i := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			s := strconv.Itoa(grid[row][col])
			solvedArray[i] = s
			i++
		}
	}
}

func isValid(board *[9][9]int) bool {

	for row := 0; row < 9; row++ {
		checker := [9]int{}
		for col := 0; col < 9; col++ {
			checker[col] = board[row][col]
			if col == 8 {
				if hasDuplicates(checker) {
					return false
				}
			}
		}
	}

	for row := 0; row < 9; row++ {
		checker := [9]int{}
		for col := 0; col < 9; col++ {
			checker[col] = board[col][row]
			if col == 8 {
				if hasDuplicates(checker) {
					return false
				}
			}
		}
	}

	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			loc := 0
			checker := [9]int{}
			for row := i; row < i+3; row++ {
				for col := j; col < j+3; col++ {
					checker[col+row+loc-j-i] = board[row][col]
				}
				loc = loc + 2
			}
			if hasDuplicates(checker) {
				return false
			}
		}
	}
	return true
}

func hasDuplicates(checker [9]int) bool {
	for i := 0; i < len(checker); i++ {
		for j := i + 1; j < len(checker); j++ {
			if checker[i] != 0 && checker[i] == checker[j] {
				return true
			}
		}
	}
	return false
}

func backtrack(board *[9][9]int) bool {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if board[row][col] == 0 {
				for i := 1; i <= 9; i++ {
					board[row][col] = i
					if isValid(board) {
						if backtrack(board) {
							return true
						}
						board[row][col] = 0
					} else {
						board[row][col] = 0
					}
				}
				return false
			}
		}
	}
	return true
}

func registerFuncs() {
	js.Global().Set("checkSudoku", js.FuncOf(checkSudoku))
	js.Global().Set("clearBoard", js.FuncOf(clearBoard))
}

func main() {
	document = js.Global().Get("document")
	c := make(chan bool)
	registerFuncs()
	<-c

}
