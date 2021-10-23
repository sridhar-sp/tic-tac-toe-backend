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

type PlayStatus string

const (
	YET_TO_START PlayStatus = "Yet to start"
	PLAYING      PlayStatus = "Playing"
	DRAW         PlayStatus = "Draw"
	PLAYER_X_WON PlayStatus = "Player x won"
	PLAYER_O_WON PlayStatus = "Player o won"
)

type Board [][]Element

type boardRepo struct {
	boardContent Board
	rowSize      int
	columnSize   int
	totalMoves   int
	playStaus    PlayStatus
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
	IsPlayerXWon() bool
	IsPlayerOWon() bool
	IsDraw() bool
	IsGameFinished() bool
	PlayComputerMove()
	GetGamePlayStatus() PlayStatus
	isWonBy(row int, column int, element Element) bool
	findGamPlayeStatus(rowIndex int, columnIndex int, element Element) PlayStatus
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
		playStaus:    YET_TO_START,
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
	b.playStaus = YET_TO_START
}

func (b *boardRepo) Mark(row int, column int, element Element) {
	b.boardContent[row][column] = element
	b.totalMoves++
	b.playStaus = b.findGamPlayeStatus(row, column, element)
}

func (b *boardRepo) GetGamePlayStatus() PlayStatus {
	return b.playStaus
}

func (b *boardRepo) IsGameStarted() bool {
	return b.playStaus != YET_TO_START
}

func (b *boardRepo) IsGameFinished() bool {
	if b.playStaus == DRAW || b.playStaus == PLAYER_X_WON || b.playStaus == PLAYER_O_WON {
		return true
	}
	return false
}

func (b *boardRepo) IsPlayerXWon() bool {
	return b.playStaus == PLAYER_X_WON
}

func (b *boardRepo) IsPlayerOWon() bool {
	return b.playStaus == PLAYER_O_WON
}

func (b *boardRepo) IsDraw() bool {
	return b.playStaus == DRAW
}

func (b *boardRepo) PlayComputerMove(element Element) (int, int) {
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

func (b boardRepo) findGamPlayeStatus(rowIndex int, columnIndex int, element Element) PlayStatus {

	//Check vertical
	for i := 0; i < b.rowSize; i++ {
		if b.GetCell(i, columnIndex) != element {
			break
		}
		if i == b.rowSize-1 {
			return elementToPlayStatus(element)
		}
	}

	//Check horizontal
	for i := 0; i < b.columnSize; i++ {
		if b.GetCell(rowIndex, i) != element {
			break
		}
		if i == b.columnSize-1 {
			return elementToPlayStatus(element)
		}
	}

	//Check dignoal (left to right)
	if rowIndex == columnIndex {
		for i := 0; i < b.rowSize; i++ {
			if b.GetCell(i, i) != element {
				break
			}
			if i == b.rowSize-1 {
				return elementToPlayStatus(element)
			}
		}
	}

	//Check dignoal (right to left)
	if rowIndex+columnIndex == b.rowSize-1 {
		for r := 0; r < b.rowSize; r++ {
			c := b.rowSize - 1 - r
			if b.GetCell(r, c) != element {
				break
			}
			if r == b.rowSize-1 {
				return elementToPlayStatus(element)
			}
		}
	}

	if b.totalMoves == b.GetBoardSize() {
		return DRAW
	}

	return PLAYING
}

func elementToPlayStatus(element Element) PlayStatus {
	switch element {
	case CROSS_CELL:
		return PLAYER_X_WON
	case NOUGHT_CELL:
		return PLAYER_O_WON
	default:
		return PLAYING
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
