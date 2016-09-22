package main

import (
	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/math"
)

//StrList TODO
type StrList struct {
	gxui.AdapterBase
	strs []string
}

//AddString TODO
func (s *StrList) AddString(add string) {
	s.strs = append(s.strs, add)
	s.DataChanged(false)
}

//Remove TODO
func (s *StrList) Remove(index int) {
	s.strs = append(s.strs[:index], s.strs[index+1:]...)
	s.DataChanged(false)
}

//SetStrings TODO
func (s *StrList) SetStrings(strs []string) {
	s.strs = strs
	s.DataChanged(false)
}

//Count TODO
func (s *StrList) Count() int {
	return len(s.strs)
}

//ItemAt TODO
func (s *StrList) ItemAt(index int) gxui.AdapterItem {
	return s.strs[index]
}

//ItemIndex TODO
func (s *StrList) ItemIndex(item gxui.AdapterItem) int {
	for i, v := range s.strs {
		if v == item {
			return i
		}
	}
	return -1
}

//Create TODO
func (s *StrList) Create(th gxui.Theme, index int) gxui.Control {
	box := th.CreateLinearLayout()
	box.SetDirection(gxui.LeftToRight)
	lbl := th.CreateLabel()
	lbl.SetText(s.strs[index])
	box.AddChild(lbl)
	return box
}

//Size TODO
func (s *StrList) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}
