package catle

import (
	"strconv"
	"sort"
)

type ColumnInfo interface {
	SetHeader(string)
	GetHeader() string
	SetWidth(int)
	GetWidth() int
	SetStatus(status)
	GetStatus() status
}

type Column interface {
	ColumnInfo
	Append(string) error
	String(int) string
	SortedOrder(bool) []int
	Size() int
}

type DefaultColumnInfo struct {
	Header string
	Width  int
	Status status
}

func (ci *DefaultColumnInfo) SetHeader(h string) {
	ci.Header = h
}

func (ci *DefaultColumnInfo) GetHeader() string {
	return ci.Header
}

func (ci *DefaultColumnInfo) SetWidth(w int) {
	ci.Width = w
}
func (ci *DefaultColumnInfo) GetWidth() int {
	return ci.Width
}

func (ci *DefaultColumnInfo) SetStatus(s status) {
	ci.Status = s
}
func (ci *DefaultColumnInfo) GetStatus() status {
	return ci.Status
}

type ColumnStr struct {
	DefaultColumnInfo
	Values []string
}

func (cs *ColumnStr) String(i int) string {
	return cs.Values[i]
}

func (cs *ColumnStr) Append(val string) error {
	cs.Values = append(cs.Values, val)
	return nil
}

func (cs *ColumnStr) Size() int {
	return len(cs.Values)
}

func (cs *ColumnStr) SortedOrder(ascending bool) []int{
	ix := make([]int, len(cs.Values))
	for i := range ix {
		ix[i] = i
	}
	if ascending {
		sort.Slice(ix, func(i, j int) bool {
			return cs.Values[ix[i]] < cs.Values[ix[j]]
		})
	}else{
		sort.Slice(ix, func(i, j int) bool {
			return cs.Values[ix[i]] > cs.Values[ix[j]]
		})
	}
	return ix
}

type ColumnInt struct {
	DefaultColumnInfo
	Values []int
}

func (ci *ColumnInt) String(i int) string {
	return strconv.Itoa(ci.Values[i])
}

func (ci *ColumnInt) Append(val string) error {
	num, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	ci.Values = append(ci.Values, num)
	return nil
}

func (ci *ColumnInt) Size() int {
	return len(ci.Values)
}

func (ci *ColumnInt) SortedOrder(ascending bool) []int {
	ix := make([]int, len(ci.Values))
	for i := range ix {
		ix[i] = i
	}
	if ascending {
		sort.Slice(ix, func(i, j int) bool {
			return ci.Values[ix[i]] < ci.Values[ix[j]]
		})
	}else{
		sort.Slice(ix, func(i, j int) bool {
			return ci.Values[ix[i]] > ci.Values[ix[j]]
		})
	}
	return ix
}

func NewColumnStr(c Column) (*ColumnStr, error) {
	ci := &ColumnStr{}
	ci.SetHeader(c.GetHeader())
	ci.SetStatus(c.GetStatus())
	ci.SetWidth(c.GetWidth())
	for i := 0; i < c.Size(); i++ {
		if err := ci.Append(c.String(i)); err != nil {
			return nil, err
		}
	}
	return ci, nil
}


func NewColumnInt(c Column) (*ColumnInt, error) {
	ci := &ColumnInt{}
	ci.SetHeader(c.GetHeader())
	ci.SetStatus(c.GetStatus())
	ci.SetWidth(c.GetWidth())
	for i := 0; i < c.Size(); i++ {
		if err := ci.Append(c.String(i)); err != nil {
			return nil, err
		}
	}
	return ci, nil
}
