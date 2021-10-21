package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const REDIRECT_URL = "https://github.com/sridhar-sp/tic-tac-toe"

const ROW_SZIE = 3
const COLUMN_SZIE = 3
const BOARD_SIZE = ROW_SZIE * COLUMN_SZIE

type Element int

const EMPTY_CELL Element = 0
const CROSS_CELL Element = 1
const NOUGHT_CELL Element = 2

var board [][]Element = getEmptyBoard(ROW_SZIE, COLUMN_SZIE)

var emptyBoxSVG []byte
var crossImageSVG []byte
var noughtImageSVG []byte
var computerStartButtonSVG []byte
var restartButtonSVG []byte

func loadAllAssets() {
	log.Println("Caching all the required asset files")
	emptyBoxSVG = fileToBytes("./assets/rect_empty.svg")
	crossImageSVG = fileToBytes("./assets/rect_x.svg")
	noughtImageSVG = fileToBytes("./assets/rect_o.svg")
	computerStartButtonSVG = fileToBytes("./assets/computer_start_button.svg")
	restartButtonSVG = fileToBytes("./assets/restart_button.svg")
}

func resetBoard() {
	board = getEmptyBoard(ROW_SZIE, COLUMN_SZIE)
}

func parseAndValidateCellIndexFromRequest(req *http.Request) (int, error) {

	cellIndex, cellIndexParseError := strconv.Atoi(req.URL.Query().Get("cellIndex"))

	if cellIndexParseError != nil {
		return -1, errors.New("Failed to parse cellIndex. received value is [" + req.URL.Query().Get("cellIndex") + "]")
	}

	log.Println("cell index is", cellIndex)

	if cellIndex >= BOARD_SIZE {
		return -1, errors.New("cell index is invalid")
	}

	return cellIndex, nil
}

func writeCommonHeaders(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "image/svg+xml")
	responseWriter.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	responseWriter.Header().Set("expires", "0")
	responseWriter.Header().Set("pragma", "no-cache")
}

func writeEmptyTextResponse(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.Write([]byte(""))
}

func getRowAndColumnIndex(cellIndex int) (int, int) {
	rowIndex := cellIndex / ROW_SZIE
	colIndex := cellIndex % COLUMN_SZIE
	return rowIndex, colIndex
}

func fileToBytes(filePath string) []byte {
	log.Println("Reading ", filePath)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error parsing ", filePath)
		return make([]byte, 0)
	}
	return bytes
}

func isGameStarted(gameBoard [][]Element) bool {
	for i := 0; i < len(gameBoard); i++ {
		for j := 0; j < len(gameBoard[i]); j++ {
			if gameBoard[i][j] != EMPTY_CELL {
				return true
			}
		}
	}
	return false
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

// Todo remove this method once initial devlopment is over
func printCells(gameBoard [][]Element) {

	for i := 0; i < len(gameBoard); i++ {
		for j := 0; j < len(gameBoard[i]); j++ {
			fmt.Print(gameBoard[i][j])
			fmt.Printf(" ")
		}
		fmt.Println("")
	}
}

// Router methods starts here

func onRenderCell(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onRenderCell function called")
	writeCommonHeaders(responseWriter)

	cellIndex, cellIndexParseError := parseAndValidateCellIndexFromRequest(req)
	if cellIndexParseError != nil {
		log.Println(cellIndexParseError.Error())
		writeEmptyTextResponse(responseWriter)
		return
	}

	rowIndex, colIndex := getRowAndColumnIndex(cellIndex)

	switch cellValue := board[rowIndex][colIndex]; cellValue {
	case CROSS_CELL:
		log.Println("Render cross cell")
		responseWriter.Write(crossImageSVG)
	case NOUGHT_CELL:
		log.Println("Render nougt cell")
		responseWriter.Write(noughtImageSVG)
	default:
		log.Println("Render empty cell")
		responseWriter.Write(emptyBoxSVG)
	}
}

func onClickCell(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onClickCell function called")
	writeCommonHeaders(responseWriter)

	cellIndex, cellIndexParseError := parseAndValidateCellIndexFromRequest(req)
	if cellIndexParseError != nil {
		log.Println(cellIndexParseError.Error())
		writeEmptyTextResponse(responseWriter)
		return
	}

	rowIndex, colIndex := getRowAndColumnIndex(cellIndex)

	board[rowIndex][colIndex] = CROSS_CELL

	http.Redirect(responseWriter, req, REDIRECT_URL, http.StatusMovedPermanently)
}

func onRenderPlayControl(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onRenderPlayControl function called")
	writeCommonHeaders(responseWriter)

	if isGameStarted(board) {
		responseWriter.Write(restartButtonSVG)
	} else {
		responseWriter.Write(computerStartButtonSVG)
	}
}

func onPlayControlClick(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onPlayControlClick function called")
	writeCommonHeaders(responseWriter)

	if isGameStarted(board) {
		// Restart button clicked
		log.Println("Reset the board")
		resetBoard()
	} else {
		// 'Let computer play first button' clicked
		log.Println("Computer making a move")
		board[1][1] = NOUGHT_CELL // Todo make a random selection
	}

	http.Redirect(responseWriter, req, REDIRECT_URL, http.StatusMovedPermanently)
}

// Router methods ends here

func main() {
	log.Println("Starting tic-tac-toe service")

	loadAllAssets()

	printCells(board)

	http.Handle("/renderCell", http.HandlerFunc(onRenderCell))
	http.Handle("/clickCell", http.HandlerFunc(onClickCell))
	http.Handle("/renderPlayControls", http.HandlerFunc(onRenderPlayControl))
	http.Handle("/clickPlayControls", http.HandlerFunc(onPlayControlClick))

	http.ListenAndServe(":2021", nil)
}
