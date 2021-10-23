package board

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Element int

const (
	EMPTY_CELL  Element = 0
	CROSS_CELL  Element = 1
	NOUGHT_CELL Element = 2
)

type Board [][]Element

type boardRepo struct {
	boardContent Board
	rowSize      int
	columnSize   int
	totalMoves   int
	randomizer   *rand.Rand
}

type BoardManger interface {
	New() *boardRepo
	GetBoardSize() int
	GetCell(row int, column int) Element
	IsEmptyCell(row int, column int) bool
	GetBoard() Board
	ResetBoard()
	Mark(row int, column int, element Element)
	IsGameStarted() bool
	PlayComputerMove()
	PrintBoard()
	getEmptyCellIndexes() []Element
	getFlatIndexFromRowAndColumn() int
	GetRowAndColumnFromFlatIndex(cellIndex int) (int, int)
}

func New(rowSize int, columnSize int) *boardRepo {
	return &boardRepo{
		boardContent: getEmptyBoard(rowSize, columnSize),
		rowSize:      rowSize,
		columnSize:   columnSize,
		totalMoves:   0,
		randomizer:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (b boardRepo) GetBoardSize() int {
	return b.rowSize * b.columnSize
}

func (b boardRepo) GetCell(row int, column int) Element {
	return b.boardContent[row][column]
}

func (b boardRepo) IsEmptyCell(row int, column int) bool {
	return b.GetCell(row, column) == EMPTY_CELL
}

func (b *boardRepo) GetBoard() Board {
	return b.boardContent
}

func (b *boardRepo) ResetBoard() {
	b.boardContent = getEmptyBoard(b.rowSize, b.columnSize)
	b.totalMoves = 0
}

func (b *boardRepo) Mark(row int, column int, element Element) {
	b.boardContent[row][column] = element
	b.totalMoves++
}

func (b *boardRepo) IsGameStarted() bool {
	for row := 0; row < b.rowSize; row++ {
		for column := 0; column < b.columnSize; column++ {
			if b.GetCell(row, column) != EMPTY_CELL {
				return true
			}
		}
	}
	return false
}

func (b boardRepo) PlayComputerMove(element Element) (int, int) {
	emptyCellIndexs := b.getEmptyCellIndexes()
	availableSpaces := len(emptyCellIndexs)

	if availableSpaces == 0 {
		log.Println("No empty space for computer to make a move")
		return -1, -1
	}

	indexToSelect := b.randomizer.Intn(availableSpaces)

	rowIndex, colIndex := b.GetRowAndColumnFromFlatIndex(emptyCellIndexs[indexToSelect])
	b.Mark(rowIndex, colIndex, element)
	return rowIndex, colIndex
}

func (b boardRepo) GetRowAndColumnFromFlatIndex(cellIndex int) (int, int) {
	rowIndex := cellIndex / b.rowSize
	colIndex := cellIndex % b.columnSize
	return rowIndex, colIndex
}

func (b boardRepo) getEmptyCellIndexes() []int {
	emptyIndexes := make([]int, 0, b.GetBoardSize())

	for row := 0; row < b.rowSize; row++ {
		for column := 0; column < b.columnSize; column++ {
			if b.GetCell(row, column) == EMPTY_CELL {
				emptyIndexes = append(emptyIndexes, b.getFlatIndexFromRowAndColumn(row, column))
			}
		}
	}

	return emptyIndexes
}

func (b boardRepo) getFlatIndexFromRowAndColumn(row int, column int) int {
	return row*b.rowSize + column
}

func (b boardRepo) PrintBoard() {

	for rowIndex := 0; rowIndex < b.rowSize; rowIndex++ {
		for columnIndex := 0; columnIndex < b.columnSize; columnIndex++ {
			fmt.Print(b.boardContent[rowIndex][columnIndex])
			fmt.Printf(" ")
		}
		fmt.Println("")
	}
}

func getEmptyBoard(rowSize int, columSize int) [][]Element {
	var boardData = make([][]Element, rowSize)

	for row := 0; row < rowSize; row++ {
		boardData[row] = make([]Element, columSize)
		for column := 0; column < columSize; column++ {
			boardData[row][column] = EMPTY_CELL
		}
	}

	return boardData
}
