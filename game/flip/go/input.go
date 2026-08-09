package flip

import "github.com/hajimehoshi/ebiten/v2"

type ClickOn int8

const (
	ClickOn_None   ClickOn = 0 // 未点击或点击位置无意义
	ClickOn_Button         = 1 // 点击按钮
	ClickOn_Card           = 2 // 点击卡片
)

type Input struct {
	clickOn        ClickOn
	clickCardIndex int // 仅点击卡片时生效，表示点击的卡片序号
}

func (ins *Input) Update(state GameState) {
	ins.clickOn = ClickOn_None
	ins.clickCardIndex = -1

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		clickX, clickY := ebiten.CursorPosition()

		switch state {
		case GameState_Prepared, GameState_End:
			if isClickOnButton(clickX, clickY) {
				ins.clickOn = ClickOn_Button
			}
		case GameState_Playing:
			index := getClickCardIndex(clickX, clickY)
			if index >= 0 {
				ins.clickOn = ClickOn_Card
				ins.clickCardIndex = index
			}
		}
	}
}

func isClickOnButton(clickPosX int, clickPosY int) bool {
	return (border+buttonMargin < clickPosX && clickPosX < border+buttonMargin+buttonWidth) &&
		(border+ButtonBoardOffsetHeight < clickPosY && clickPosY < border+ButtonBoardOffsetHeight+buttonHeight)
}

func getClickCardIndex(clickPosX int, clickPosY int) int {
	if !(border+CardBoardOffsetWidth < clickPosX && clickPosX < border+CardBoardOffsetWidth+CardBoardWidth) ||
		!(border+CardBoardOffsetHeight < clickPosY && clickPosY < border+CardBoardOffsetHeight+CardBoardHeight) {
		return -1 // click outside of cardboard
	}

	col := 0
	for i := 1; i < 5; i++ {
		if border+CardBoardOffsetWidth+i*cardMargin+(i-1)*cardWidth < clickPosX &&
			clickPosX < border+CardBoardOffsetWidth+i*cardMargin+i*cardWidth {
			col = i
			break
		}
	}
	if col == 0 {
		return -1
	}

	row := 0
	for i := 1; i < 5; i++ {
		if border+CardBoardOffsetHeight+i*cardMargin+(i-1)*cardHeight < clickPosY &&
			clickPosY < border+CardBoardOffsetHeight+i*cardMargin+i*cardHeight {
			row = i
			break
		}
	}
	if row == 0 {
		return -1
	}

	return (row-1)*4 + (col - 1)
}
