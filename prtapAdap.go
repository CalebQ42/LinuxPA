package main

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"path"
	"sort"
	"strings"

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
	box.SetPadding(math.CreateSpacing(2))
	box.SetDirection(gxui.LeftToRight)
	box.SetVerticalAlignment(gxui.AlignMiddle)
	dir := path.Dir(p.apps[index].ex)
	if fold, err := os.Open(dir + "/App/AppInfo"); err == nil {
		var pics []string
		fi, _ := fold.Readdirnames(-1)
		for _, v := range fi {
			if strings.HasPrefix(v, "appicon_") && strings.HasSuffix(v, ".png") {
				pics = append(pics, v)
			}
		}
		if len(pics) > 0 {
			ind := sort.SearchStrings(pics, "appicon_128.png")
			if ind == len(pics) {
				ind = len(pics) - 1
			}
			imgfi, _ := os.Open(dir + "/App/AppInfo/" + pics[ind])
			img, _, err := image.Decode(imgfi)
			if err == nil {
				rgba := image.NewRGBA(img.Bounds())
				draw.Draw(rgba, img.Bounds(), img, image.ZP, draw.Src)
				tex := dr.CreateTexture(rgba, 1)
				icon := th.CreateImage()
				icon.SetExplicitSize(math.Size{H: 32, W: 32})
				icon.SetTexture(tex)
				box.AddChild(icon)
			}
		}
	} else if fi, err := os.Open(dir + "/appicon.png"); err == nil {
		img, _, err := image.Decode(fi)
		if err == nil {
			rgba := image.NewRGBA(img.Bounds())
			draw.Draw(rgba, img.Bounds(), img, image.ZP, draw.Src)
			tex := dr.CreateTexture(rgba, 1)
			icon := th.CreateImage()
			icon.SetExplicitSize(math.Size{H: 32, W: 32})
			icon.SetTexture(tex)
			box.AddChild(icon)
		}
	} else {
		//Creating empty Image so that names line up
		icon := th.CreateImage()
		icon.SetExplicitSize(math.Size{H: 32, W: 32})
		box.AddChild(icon)
	}
	lbl := th.CreateLabel()
	lbl.SetText(p.apps[index].name)
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
	return math.Size{W: math.MaxSize.W, H: 36}
}
