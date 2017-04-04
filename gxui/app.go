package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/nelsam/gxui"
	"github.com/nelsam/gxui/math"
)

type app struct {
	name   string
	cat    string
	appimg []string
	lin    []string
	ex     []string
	icon   gxui.Texture
	dir    string
	ini    *os.File
}

type appExNode struct {
	ap    app
	exInd int
}

func (a *appExNode) launch() {
	if wine {
		var cmd *exec.Cmd
		if !contains(a.ap.lin, a.ap.ex[a.exInd]) {
			cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; wine \""+a.ap.ex[a.exInd]+"\"")
		} else {
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[a.exInd]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[a.exInd]+"\"")
			}
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}
	var cmd *exec.Cmd
	if comEnbld {
		cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[a.exInd]+"\"")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[a.exInd]+"\"")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}

func (a *appExNode) Count() int {
	return 0
}

func (a *appExNode) NodeAt(int) gxui.TreeNode {
	return nil
}

func (a *appExNode) ItemIndex(gxui.AdapterItem) int {
	return -1
}

func (a *appExNode) Item() gxui.AdapterItem {
	if wine {
		return a.ap.ex[a.exInd]
	}
	return a.ap.lin[a.exInd]
}

func (a *appExNode) Create(the gxui.Theme) gxui.Control {
	box := the.CreateLinearLayout()
	box.SetDirection(gxui.LeftToRight)
	box.SetVerticalAlignment(gxui.AlignMiddle)
	img := the.CreateImage()
	img.SetTexture(a.ap.icon)
	img.SetExplicitSize(math.Size{H: 32, W: 32})
	lbl := the.CreateLabel()
	lbl.SetText(a.ap.ex[a.exInd])
	box.AddChild(img)
	box.AddChild(lbl)
	box.OnDoubleClick(func(gxui.MouseEvent) {
		a.launch()
	})
	return box
}

type appNode struct {
	ap app
}

func (a *appNode) launch() {
	if len(a.ap.ex) == 1 {
		if wine {
			var cmd *exec.Cmd
			if !contains(a.ap.lin, a.ap.ex[0]) {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; wine \""+a.ap.ex[0]+"\"")
			} else {
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[0]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[0]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			var cmd *exec.Cmd
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[0]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.ex[0]+"\"")
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}
	} else {
		if wine {
			var cmd *exec.Cmd
			if len(a.ap.lin) == 0 {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; wine \""+a.ap.ex[0]+"\"")
			} else {
				var ind int
				for i, v := range a.ap.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.lin[ind]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			if len(a.ap.lin) != 0 {
				var ind int
				for i, v := range a.ap.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				var cmd *exec.Cmd
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.ap.dir+"\"; \"./"+a.ap.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.ap.dir+"\"; \"./"+a.ap.lin[ind]+"\"")
				}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
			}
		}
	}
}

func (a *appNode) Count() int {
	if wine {
		if len(a.ap.ex) > 1 {
			return len(a.ap.ex)
		}
		return 0
	}
	if len(a.ap.lin) > 1 {
		return len(a.ap.lin)
	}
	return 0
}

func (a *appNode) NodeAt(i int) gxui.TreeNode {
	return &appExNode{ap: a.ap, exInd: i}
}

func (a *appNode) ItemIndex(item gxui.AdapterItem) int {
	if wine {
		for i, v := range a.ap.ex {
			if v == item {
				return i
			}
		}
	} else {
		for i, v := range a.ap.lin {
			if v == item {
				return i
			}
		}
	}
	return -1
}

func (a *appNode) Item() gxui.AdapterItem {
	return a.ap.name
}

func (a *appNode) Create(the gxui.Theme) gxui.Control {
	box := the.CreateLinearLayout()
	box.SetDirection(gxui.LeftToRight)
	box.SetPadding(math.CreateSpacing(2))
	box.SetVerticalAlignment(gxui.AlignMiddle)
	img := the.CreateImage()
	if a.ap.icon != nil {
		img.SetTexture(a.ap.icon)
	}
	img.SetExplicitSize(math.Size{H: 32, W: 32})
	lbl := the.CreateLabel()
	lbl.SetText(a.ap.name)
	box.AddChild(img)
	box.AddChild(lbl)
	box.OnDoubleClick(func(gxui.MouseEvent) {
		a.launch()
	})
	return box
}

type catAdap struct {
	gxui.AdapterBase
	cat string
}

func (a *catAdap) setCat(cat string) {
	a.cat = cat
	a.DataChanged(false)
}

func (a *catAdap) refresh() {
	a.DataChanged(false)
}

func (a *catAdap) Count() int {
	if wine {
		return len(master[a.cat])
	}
	return len(linmaster[a.cat])
}

func (a *catAdap) NodeAt(i int) gxui.TreeNode {
	if wine {
		return &appNode{ap: master[a.cat][i]}
	}
	return &appNode{ap: linmaster[a.cat][i]}
}

func (a *catAdap) Size(gxui.Theme) math.Size {
	return math.Size{H: 34, W: math.MaxSize.W}
}

func (a *catAdap) ItemIndex(item gxui.AdapterItem) int {
	if wine {
		for i, v := range master[a.cat] {
			if v.name == item {
				return i
			}
		}
	} else {
		for i, v := range linmaster[a.cat] {
			if v.name == item {
				return i
			}
		}
	}
	return -1
}
