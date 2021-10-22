package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/sridhar-sp/tic-tac-toe-backend/src/assets"
)

const REDIRECT_URL = "https://github.com/sridhar-sp/tic-tac-toe"

const ROW_SIZE = 3
const COLUMN_SIZE = 3
const BOARD_SIZE = ROW_SIZE * COLUMN_SIZE

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type Element int

const EMPTY_CELL Element = 0
const CROSS_CELL Element = 1
const NOUGHT_CELL Element = 2

var board [][]Element = getEmptyBoard(ROW_SIZE, COLUMN_SIZE)

var activities = []string{"* Game is yet to begin"}

func resetBoard() {
	board = getEmptyBoard(ROW_SIZE, COLUMN_SIZE)
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
	rowIndex := cellIndex / ROW_SIZE
	colIndex := cellIndex % COLUMN_SIZE
	return rowIndex, colIndex
}

func getFlatIndexFromRowAndColumn(row int, column int, rowSize int, columnSize int) int {
	return row*rowSize + column
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

func getEmptyCellIndexes(gameBoard [][]Element) []int {
	emptyIndexes := make([]int, 0, BOARD_SIZE)

	rowSize := len(gameBoard)
	for row := 0; row < rowSize; row++ {
		columnSize := len(gameBoard[row])
		for column := 0; column < columnSize; column++ {
			if gameBoard[row][column] == EMPTY_CELL {
				emptyIndexes = append(emptyIndexes, getFlatIndexFromRowAndColumn(row, column, rowSize, columnSize))
			}
		}
	}

	return emptyIndexes
}

func playComputerMove(gameBoard [][]Element) (int, int) {
	emptyCellIndexs := getEmptyCellIndexes(gameBoard)
	availableSpaces := len(emptyCellIndexs)
	if availableSpaces == 0 {
		log.Println("No empty space for computer to make a move")
		return -1, -1
	}

	indexToSelect := random.Intn(availableSpaces)

	rowIndex, colIndex := getRowAndColumnIndex(emptyCellIndexs[indexToSelect])
	gameBoard[rowIndex][colIndex] = NOUGHT_CELL
	return rowIndex, colIndex
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

func getPort() string {
	var port = "3000"
	var portFromConfig, isConfigPresent = os.LookupEnv("PORT")
	if isConfigPresent {
		port = portFromConfig
	}
	return port
}

func writeTextResponseAsImage(textLines []string, responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "image/svg+xml")
	drawTextLines(textLines, 18, 8, "#58A6FF", "start", responseWriter)
}

func drawTextLines(textLines []string, fontSize int, lineSpace int, textColor string, textAlign string, writer io.Writer) {
	var longestTextLength = -1
	for _, text := range textLines {
		currentTextLength := len(text)
		if currentTextLength > longestTextLength {
			longestTextLength = currentTextLength
		}
	}

	width := int(float32(longestTextLength) * (float32(fontSize) / float32(2)))
	height := fontSize*len(textLines) + (len(textLines) * lineSpace)
	canvas := svg.New(writer)
	canvas.Start(width, height)

	// canvas.Rect(0, 0, width, height, "fill:#238636")

	var textStartX = 0
	switch textAlign {
	case "start":
		textStartX = 0
	case "middle":
		textStartX = width / 2
	case "end":
		textStartX = width
	}

	style := fmt.Sprintf("font-size:%dpx;fill:%s;text-anchor:%s", fontSize, textColor, textAlign)
	var textStartY = fontSize - (fontSize / 5)
	for _, text := range textLines {
		canvas.Text(textStartX, textStartY, text, style)
		textStartY = textStartY + fontSize + lineSpace
	}

	canvas.End()
}

func updateActivities(text string) {
	activities = append(activities, text)
}

func clearActivities() {
	activities = []string{}
}

func getActivities() []string {
	return activities
}

// Router methods starts here

func onHome(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onHome function called")

	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.Write([]byte("Welcome to tic-tac-toe. application is running at port " + getPort()))
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
		responseWriter.Write(assetManager.GetSvg(assets.CrossImage))
	case NOUGHT_CELL:
		log.Println("Render nougt cell")
		responseWriter.Write(assetManager.GetSvg(assets.NoughtImage))
	default:
		log.Println("Render empty cell")
		responseWriter.Write(assetManager.GetSvg(assets.EmptyBox))
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
	computerMoveRowIndex, computerMoveColIndex := playComputerMove(board)

	updateActivities(fmt.Sprintf("* Player X clicked cell [%d,%d]", rowIndex, colIndex))
	updateActivities(fmt.Sprintf("* Player O clicked cell [%d,%d]", computerMoveRowIndex, computerMoveColIndex))

	http.Redirect(responseWriter, req, REDIRECT_URL, http.StatusMovedPermanently)
}

func onRenderPlayControl(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onRenderPlayControl function called")
	writeCommonHeaders(responseWriter)

	if isGameStarted(board) {
		responseWriter.Write(assetManager.GetSvg(assets.RestartButton))
	} else {
		responseWriter.Write(assetManager.GetSvg(assets.ComputerStartButton))
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

func renderActivities(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("renderActivities function called")
	writeCommonHeaders(responseWriter)

	writeTextResponseAsImage(activities, responseWriter)
}

// Router methods ends here

var assetManager = assets.New()

func main() {
	port := getPort()
	log.Println("Starting tic-tac-toe service at port ", port)

	http.Handle("/", http.HandlerFunc(onHome))
	http.Handle("/renderCell", http.HandlerFunc(onRenderCell))
	http.Handle("/clickCell", http.HandlerFunc(onClickCell))
	http.Handle("/renderPlayControls", http.HandlerFunc(onRenderPlayControl))
	http.Handle("/clickPlayControls", http.HandlerFunc(onPlayControlClick))
	http.Handle("/renderActivities", http.HandlerFunc(renderActivities))

	http.ListenAndServe(":"+port, nil)
}
