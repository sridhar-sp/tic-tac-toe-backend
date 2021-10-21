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

var board = [3][3]Element{
	{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
	{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
	{EMPTY_CELL, EMPTY_CELL, EMPTY_CELL},
}

var emptyBoxSVG []byte
var crossImageSVG []byte
var noughtImageSVG []byte

func loadAllAssets() {
	log.Println("Caching all the required asset files")
	emptyBoxSVG = fileToBytes("./assets/rect_empty.svg")
	crossImageSVG = fileToBytes("./assets/rect_x.svg")
	noughtImageSVG = fileToBytes("./assets/rect_o.svg")
}

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

func printCells() {

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(board[i][j])
			fmt.Printf(" ")
		}
		fmt.Println("")
	}
}

func main() {
	log.Println("Starting tic-tac-toe service")

	loadAllAssets()

	printCells()

	http.Handle("/renderCell", http.HandlerFunc(onRenderCell))
	http.Handle("/clickCell", http.HandlerFunc(onClickCell))

	http.ListenAndServe(":2021", nil)
}
