package main

import (
	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/math"
)

type prtapAdap struct {
	gxui.AdapterBase
	apps []prtap
}

func (p *prtapAdap) SetApps(apps []prtap) {
	p.apps = apps
	p.DataChanged(false)
}

func (p *prtapAdap) Count() int {
	return len(p.apps)
}

func (p *prtapAdap) Create(th gxui.Theme, index int) gxui.Control {
	box := th.CreateLinearLayout()
	box.SetDirection(gxui.LeftToRight)
	//add image support
	// pic := th.CreateImage()
	// dr.CreateTexture()
	lbl := th.CreateLabel()
	lbl.SetText(p.apps[index].name)
	// box.AddChild(pic)
	box.AddChild(lbl)
	return box
}

func (p *prtapAdap) ItemAt(index int) gxui.AdapterItem {
	return p.apps[index]
}

func (p *prtapAdap) ItemIndex(item gxui.AdapterItem) int {
	it, ok := item.(prtap)
	if !ok {
		return -1
	}
	for i, v := range p.apps {
		if v == it {
			return i
		}
	}
	return -1
}

func (p *prtapAdap) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}
