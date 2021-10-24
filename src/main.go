package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
	"github.com/sridhar-sp/tic-tac-toe-backend/src/assets"
	"github.com/sridhar-sp/tic-tac-toe-backend/src/board"
)

const REDIRECT_URL = "https://github.com/sridhar-sp/tic-tac-toe"

var activities = []string{}

func parseAndValidateCellIndexFromRequest(req *http.Request) (int, error) {

	cellIndex, cellIndexParseError := strconv.Atoi(req.URL.Query().Get("cellIndex"))

	if cellIndexParseError != nil {
		return -1, errors.New("Failed to parse cellIndex. received value is [" + req.URL.Query().Get("cellIndex") + "]")
	}

	log.Println("cell index is", cellIndex)

	if cellIndex >= boardManager.GetBoardSize() {
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
	drawTextLines(textLines, 18, 8, "#C9D1D9", "#58A6FF", "start", responseWriter)
}

func drawTextLines(textLines []string, fontSize int, lineSpace int, textColor string, lastRecentTextColor string, textAlign string, writer io.Writer) {
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

	oldTextStyle := fmt.Sprintf("font-size:%dpx;fill:%s;text-anchor:%s", fontSize, textColor, textAlign)
	lastRecentTextStyle := fmt.Sprintf("font-size:%dpx;fill:%s;text-anchor:%s", fontSize, lastRecentTextColor, textAlign)

	lastIndex := len(textLines) - 1
	var textStartY = fontSize
	for index, text := range textLines {
		if index == lastIndex {
			canvas.Text(textStartX, textStartY, text, lastRecentTextStyle)
		} else {
			canvas.Text(textStartX, textStartY, text, oldTextStyle)
		}

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
	if len(activities) == 0 {
		return []string{"* Game is yet to begin"}
	}
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

	rowIndex, colIndex := boardManager.GetRowAndColumnFromFlatIndex(cellIndex)
	isCellSelected := boardManager.IsCellSelectedByWinner(rowIndex, colIndex)

	switch cellValue := boardManager.GetCell(rowIndex, colIndex); cellValue {
	case board.CROSS_CELL:
		log.Println("Render cross cell")
		if isCellSelected {
			responseWriter.Write(assetManager.GetSvg(assets.SelectedCrossImage))
		} else {
			responseWriter.Write(assetManager.GetSvg(assets.CrossImage))
		}
	case board.NOUGHT_CELL:
		log.Println("Render nougt cell")
		if isCellSelected {
			responseWriter.Write(assetManager.GetSvg(assets.SelectedNoughtImage))
		} else {
			responseWriter.Write(assetManager.GetSvg(assets.NoughtImage))
		}
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

	if !boardManager.IsGameFinished() {
		rowIndex, colIndex := boardManager.GetRowAndColumnFromFlatIndex(cellIndex)
		if boardManager.IsEmptyCell(rowIndex, colIndex) {
			boardManager.Mark(rowIndex, colIndex, board.CROSS_CELL)
			updateActivities(fmt.Sprintf("* Player 'X' clicked cell [%d , %d]", rowIndex, colIndex))

			if boardManager.IsPlayerXWon() {
				updateActivities("* Player 'X' won")
			} else if boardManager.IsDraw() {
				updateActivities("* Game draw")
			} else {
				computerMoveRowIndex, computerMoveColIndex := boardManager.PlayComputerMove(board.NOUGHT_CELL)
				updateActivities(fmt.Sprintf("* Player 'O' clicked cell [%d , %d]", computerMoveRowIndex, computerMoveColIndex))
				if boardManager.IsPlayerOWon() {
					updateActivities("* Player 'O' won")
				} else if boardManager.IsDraw() {
					updateActivities("* Game draw")
				}
			}
		} else {
			updateActivities(fmt.Sprintf("* Player 'X' can't select cell [%d , %d], the cell is already selected", rowIndex, colIndex))
		}
	} else {
		updateActivities("* Game over. Please click 'restart' to play the game again")
	}

	http.Redirect(responseWriter, req, REDIRECT_URL, http.StatusMovedPermanently)
}

func onRenderPlayControl(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onRenderPlayControl function called")
	writeCommonHeaders(responseWriter)

	if boardManager.IsGameStarted() {
		responseWriter.Write(assetManager.GetSvg(assets.RestartButton))
	} else {
		responseWriter.Write(assetManager.GetSvg(assets.ComputerStartButton))
	}
}

func onPlayControlClick(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("onPlayControlClick function called")
	writeCommonHeaders(responseWriter)

	if boardManager.IsGameStarted() {
		log.Println("Reset the board")
		boardManager.ResetBoard()
		clearActivities()
	} else {
		log.Println("Computer making a move")
		computerMoveRowIndex, computerMoveColIndex := boardManager.PlayComputerMove(board.NOUGHT_CELL)
		updateActivities(fmt.Sprintf("* Player 'O' clicked cell [%d , %d]", computerMoveRowIndex, computerMoveColIndex))
	}

	http.Redirect(responseWriter, req, REDIRECT_URL, http.StatusMovedPermanently)
}

func renderActivities(responseWriter http.ResponseWriter, req *http.Request) {
	log.Println("renderActivities function called")
	writeCommonHeaders(responseWriter)

	writeTextResponseAsImage(getActivities(), responseWriter)
}

// Router methods ends here

const ROW_SIZE = 3
const COLUMN_SIZE = 3

var assetManager = assets.New()
var boardManager = board.New(ROW_SIZE, COLUMN_SIZE)

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
