package catle

import (
	"encoding/csv"
	"strings"
	"unicode/utf8"
	"io"
)

type ColumnFeeder struct {
	*csv.Reader // csv reader
	columns []Column
}

func NewColumnFeeder(reader *csv.Reader) *ColumnFeeder {
	return &ColumnFeeder{reader, nil}
}

func (cr *ColumnFeeder) ColNum() int { return len(cr.columns) }
func (cr *ColumnFeeder) RowNum() int { return cr.columns[0].Size() }

func (cr *ColumnFeeder) Init(withoutHeader bool) error {
	header, err := cr.Read()
	if err != nil {
		return err
	}

	cr.columns = make([]Column, len(header))
	for i := 0; i < len(header); i++ {
		cr.columns[i] = &ColumnStr{}
	}
	for i := range header {
		header[i] = strings.TrimSpace(header[i])
		if len(header[i]) == 0 {
			header[i] = " "
		}
		if withoutHeader {
			cr.columns[i].SetHeader("")
			cr.columns[i].Append(header[i])
		} else {
			cr.columns[i].SetHeader(header[i])
		}
		cr.columns[i].SetWidth(utf8.RuneCountInString(header[i]))
	}
	return nil
}

func (cr *ColumnFeeder) FeedLine(num int) (int, error) {
	readNum := 0
	for i := 0; i < num; i++ {
		line, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return readNum, err
		}
		for j, cell := range line {
			cell = strings.TrimSpace(cell)
			if len(cell) == 0 {
				cell = " "
			}
			cr.columns[j].Append(cell)
			w := utf8.RuneCountInString(cell)
			if cr.columns[j].GetWidth() < w {
				cr.columns[j].SetWidth(w)
			}
		}
		readNum += 1
	}
	return readNum, nil
}

func (cr *ColumnFeeder) FeedAll() error {
	readNum := 0
	line, err := cr.Read()
	for err == nil {
		for i, cell := range line {
			cell = strings.TrimSpace(cell)
			if len(cell) == 0 {
				cell = " "
			}
			cr.columns[i].Append(cell)
			w := utf8.RuneCountInString(cell)
			if cr.columns[i].GetWidth() < w {
				cr.columns[i].SetWidth(w)
			}
		}
		readNum += 1
		line, err = cr.Read()
	}
	if err == io.EOF {
		return nil
	}
	return err
}
