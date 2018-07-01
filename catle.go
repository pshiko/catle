package catle

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"encoding/csv"
	"log"
)

type status int

const (
	NORMAL   = iota
	SELECTED
	HIDDEN
)

type Catle struct {
	*ColumnFeeder
	delim        rune // delimiter for cell
	colBegIx     int  // begin column index for viewing
	colEndIx     int  // end column index for viewing
	curPosX      int  // cursor pos
	curPosY      int  // cursor pos
	columnPos    []int
	headerHeight int
	winWidth     int   // window width
	winHeight    int   // window height
	drawOrder    []int // drawOrder
}

func NewCatle(sc *csv.Reader, withoutHeader bool) (*Catle, error) {
	if err := termbox.Init(); err != nil {
		fmt.Errorf("error termbox.Init()")
		panic(err)
	}
	cf := NewColumnFeeder(sc)
	cf.Init(withoutHeader)
	cf.columns[0].SetStatus(SELECTED)
	width, height := termbox.Size()
	headerHeight := 1
	if withoutHeader {
		headerHeight = 0
	}
	catle := &Catle{
		ColumnFeeder: cf,
		delim:        '|',
		colBegIx:     0,
		colEndIx:     0,
		winWidth:     width,
		winHeight:    height,
		curPosX:      0,
		curPosY:      0,
		columnPos:    make([]int, cf.ColNum()),
		headerHeight: headerHeight,
	}
	return catle, nil
}

func (ca *Catle) updateWinSize() {
	ca.winWidth, ca.winHeight = termbox.Size()
}

func (ca *Catle) updatePosX(move int) {
	ca.curPosX += move
	if ca.curPosX < 0 {
		ca.curPosX = 0
	}
	if ca.curPosX >= ca.ColNum() {
		ca.curPosX = ca.ColNum() - 1
	}
}

func (ca *Catle) updatePosY(move int) {
	ca.curPosY += move
	if ca.curPosY < 0 {
		ca.curPosY = 0
	}
	if ca.curPosY >= ca.RowNum() {
		ca.curPosY = ca.RowNum() - 1
	}
}

func (ca *Catle) calcColPosition() int {
	maxCol := ca.colBegIx + 1
	ca.columnPos[ca.colBegIx] = 0
	for i := ca.colBegIx + 1; i < ca.ColNum(); i++ {
		if ca.columns[i-1].GetStatus() == HIDDEN {
			ca.columnPos[i] = ca.columnPos[i-1] + 1 + 1
		} else {
			ca.columnPos[i] = ca.columnPos[i-1] + ca.columns[i-1].GetWidth() + 1
		}
		if ca.columnPos[i] < ca.winWidth {
			maxCol = i
		}
	}
	return maxCol + 1
}

func (ca *Catle) PrintTable() error {
	ca.Clear()
	begRow := ca.curPosY
	endRow := begRow + ca.winHeight - ca.headerHeight
	if endRow > ca.RowNum() {
		_, err := ca.FeedLine(endRow - ca.RowNum())
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
	}
	endRow = ca.RowNum() - ca.headerHeight

	ca.colEndIx = ca.calcColPosition()
	for ca.curPosX >= ca.colEndIx {
		ca.colBegIx += 1
		ca.colEndIx = ca.calcColPosition()
	}

	if ca.curPosX < ca.colBegIx {
		ca.colBegIx = ca.curPosX
	}

	if ca.headerHeight > 0 {
		ca.PrintHeader()
	}
	for i := ca.colBegIx; i < ca.colEndIx; i++ {
		if i > 0 {
			ca.PrintColDelim(endRow-begRow+ca.headerHeight, ca.columnPos[i]-1)
		}
		ca.PrintColumn(i, begRow, endRow)
	}
	ca.Flush()
	return nil
}

func (ca *Catle) PrintHeader() {
	for i := ca.colBegIx; i < ca.colEndIx; i++ {
		header := ca.columns[i].GetHeader()
		fg := termbox.ColorCyan
		bg := termbox.ColorDefault

		if _, ok := ca.columns[i].(*ColumnStr); !ok {
			fg = termbox.ColorGreen
		}

		switch ca.columns[i].GetStatus() {
		case HIDDEN:
			ca.PrintCell(header[0:1], ca.columnPos[i], 0, fg, bg)
		default:
			ca.PrintCell(header, ca.columnPos[i], 0, fg, bg)
		}
	}
}

func (ca *Catle) PrintColumn(n, beg, end int) {
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if n == ca.curPosX {
		fg = termbox.ColorRed
	}

	if ca.columns[n].GetStatus() == HIDDEN {
		for i := beg; i < end; i++ {
			if len(ca.drawOrder) > 0 {
				ca.PrintCell(ca.columns[n].String(ca.drawOrder[i])[0:1], ca.columnPos[n], (i-beg)+ca.headerHeight, fg, bg)
			} else {
				ca.PrintCell(ca.columns[n].String(i)[0:1], ca.columnPos[n], (i-beg)+ca.headerHeight, fg, bg)
			}
		}
		return
	}
	for i := beg; i < end; i++ {
		if len(ca.drawOrder) > 0 {
			ca.PrintCell(ca.columns[n].String(ca.drawOrder[i]), ca.columnPos[n], (i-beg)+ca.headerHeight, fg, bg)
		} else {
			ca.PrintCell(ca.columns[n].String(i), ca.columnPos[n], (i-beg)+ca.headerHeight, fg, bg)
		}
	}
	return
}

func (ca *Catle) PrintCell(cell string, x, y int, fg termbox.Attribute, bg termbox.Attribute) {
	for i, ch := range cell {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func (ca *Catle) PrintColDelim(h int, x int) {
	for i := 0; i < h; i++ {
		termbox.SetCell(x, i, ca.delim, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (ca *Catle) PollEvent() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			ca.updateWinSize()
			ca.curPosX = ca.colBegIx
			ca.PrintTable()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeySpace:
				if ca.columns[ca.curPosX].GetStatus() != HIDDEN {
					ca.columns[ca.curPosX].SetStatus(HIDDEN)
				} else {
					ca.columns[ca.curPosX].SetStatus(NORMAL)
				}
				ca.PrintTable()
			case termbox.KeyCtrlA:
				if err := ca.FeedAll(); err != nil {
					fmt.Errorf("feeding was failed")
					panic(err)
				}
				ca.PrintTable()
			case termbox.KeyCtrlH:
				ca.updatePosX(-(ca.colEndIx - ca.colBegIx) / 2)
				ca.PrintTable()
			case termbox.KeyCtrlI:
				if ci, err := NewColumnInt(ca.columns[ca.curPosX]); err == nil {
					ca.columns[ca.curPosX] = ci
				}
				ca.PrintTable()
			case termbox.KeyCtrlK:
				ca.drawOrder = []int{}
				for i := range ca.columns {
					if _, ok := ca.columns[i].(*ColumnStr); !ok {
						cs, err := NewColumnStr(ca.columns[i])
						if err != nil {
							log.Fatalln(err)
						}
						ca.columns[i] = cs
					}
				}
				ca.PrintTable()
			case termbox.KeyCtrlL:
				ca.updatePosX((ca.colEndIx - ca.colBegIx) / 2)
				ca.PrintTable()
			case termbox.KeyCtrlN, termbox.KeyCtrlF:
				ca.updatePosY(ca.winHeight / 2)
				ca.PrintTable()
			case termbox.KeyCtrlU:
				ca.updatePosY(-ca.winHeight / 2)
				ca.PrintTable()
			case termbox.KeyArrowLeft:
				ca.updatePosX(-1)
				ca.PrintTable()
			case termbox.KeyArrowRight:
				ca.updatePosX(1)
				ca.PrintTable()
			case termbox.KeyArrowUp:
				ca.updatePosY(-1)
				ca.PrintTable()
			case termbox.KeyArrowDown:
				ca.updatePosY(1)
				ca.PrintTable()
			default:
				switch ev.Ch {
				case 'G':
					if err := ca.FeedAll(); err != nil {
						fmt.Errorf("feeding was failed")
						panic(err)
					}
					ca.updatePosY(ca.RowNum() - ca.curPosY)
					ca.PrintTable()
				case 'h':
					ca.updatePosX(-1)
					ca.PrintTable()
				case 'j':
					ca.updatePosY(1)
					ca.PrintTable()
				case 'k':
					ca.updatePosY(-1)
					ca.PrintTable()
				case 'l':
					ca.updatePosX(1)
					ca.PrintTable()
				case 'n':
					if ca.headerHeight == 1 {
						ca.headerHeight = 0
					} else {
						ca.headerHeight = 1
					}
					ca.PrintTable()
				case 's':
					if err := ca.FeedAll(); err != nil {
						fmt.Errorf("feeding was failed")
						panic(err)
					}
					ca.drawOrder = ca.columns[ca.curPosX].SortedOrder(true)
					ca.PrintTable()
				case 'S':
					if err := ca.FeedAll(); err != nil {
						fmt.Errorf("feeding was failed")
						panic(err)
					}
					ca.drawOrder = ca.columns[ca.curPosX].SortedOrder(false)
					ca.PrintTable()
				case 'q':
					return
				case 'U':
					ca.updatePosY(-ca.curPosY)
					ca.PrintTable()
				}
			}
		}
	}
}

func (ca *Catle) Flush() {
	termbox.Flush()
}

func (ca *Catle) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (ca *Catle) Close() {
	termbox.Close()
}
