package assets

import (
	"io/ioutil"
	"log"
)

type svgImageAssets struct {
	emptyBoxSVG            []byte
	crossImageSVG          []byte
	selectedCrossImageSVG  []byte
	noughtImageSVG         []byte
	selectedNoughtImageSVG []byte
	computerStartButtonSVG []byte
	restartButtonSVG       []byte
}

type AssetId int

const (
	EmptyBox AssetId = iota
	CrossImage
	SelectedCrossImage
	NoughtImage
	SelectedNoughtImage
	ComputerStartButton
	RestartButton
)

type SvgImages interface {
	New() svgImageAssets
	GetSvg(assetId AssetId) []byte
}

func New() svgImageAssets {
	log.Println("Caching all the required asset files")
	return svgImageAssets{
		emptyBoxSVG:            fileToBytes("./assets/rect_empty.svg"),
		crossImageSVG:          fileToBytes("./assets/rect_x.svg"),
		selectedCrossImageSVG:  fileToBytes("./assets/rect_x_selected.svg"),
		noughtImageSVG:         fileToBytes("./assets/rect_o.svg"),
		selectedNoughtImageSVG: fileToBytes("./assets/rect_o_selected.svg"),
		computerStartButtonSVG: fileToBytes("./assets/computer_start_button.svg"),
		restartButtonSVG:       fileToBytes("./assets/restart_button.svg"),
	}
}

func (s svgImageAssets) GetSvg(assetId AssetId) []byte {
	switch assetId {
	case EmptyBox:
		return s.emptyBoxSVG
	case CrossImage:
		return s.crossImageSVG
	case SelectedCrossImage:
		return s.selectedCrossImageSVG
	case NoughtImage:
		return s.noughtImageSVG
	case SelectedNoughtImage:
		return s.selectedNoughtImageSVG
	case ComputerStartButton:
		return s.computerStartButtonSVG
	case RestartButton:
		return s.restartButtonSVG
	default:
		log.Println("Invalid asset id")
		return []byte{}
	}
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
