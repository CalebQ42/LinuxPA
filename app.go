package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type app struct {
	name   string
	cat    string
	appimg []string
	lin    []string
	ex     []string
	icon   *gdk.Pixbuf
	dir    string
	ini    *os.File
}

func (a *app) getTreeIter(store *gtk.TreeStore) *gtk.TreeIter {
	it := store.Append(nil)
	store.SetValue(it, 0, a.icon)
	store.SetValue(it, 1, a.name)
	if len(a.ex) > 1 {
		for _, v := range a.ex {
			i := store.Append(it)
			store.SetValue(i, 1, v)
		}
	}
	return it
}

func (a *app) launch() {
	if len(a.ex) == 1 {
		if wine {
			var cmd *exec.Cmd
			if !contains(a.lin, a.ex[0]) {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
			} else {
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			var cmd *exec.Cmd
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[0]+"\"")
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		}
	} else {
		if wine {
			var cmd *exec.Cmd
			if len(a.lin) == 0 {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[0]+"\"")
			} else {
				var ind int
				for i, v := range a.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				}
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
		} else {
			if len(a.lin) != 0 {
				var ind int
				for i, v := range a.lin {
					if strings.HasSuffix(v, ".sh") {
						ind = i
						break
					}
				}
				var cmd *exec.Cmd
				if comEnbld {
					cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				} else {
					cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.lin[ind]+"\"")
				}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Start()
			}
		}
	}
}

func (a *app) launchSub(sub int) {
	if wine {
		var cmd *exec.Cmd
		if !contains(a.lin, a.ex[sub]) {
			cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; wine \""+a.ex[sub]+"\"")
		} else {
			if comEnbld {
				cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
			} else {
				cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
			}
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}
	var cmd *exec.Cmd
	if comEnbld {
		cmd = exec.Command("/bin/sh", "-c", ". PortableApps/LinuxPACom/common.sh || exit 1;cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
	} else {
		cmd = exec.Command("/bin/sh", "-c", "cd \""+a.dir+"\"; \"./"+a.ex[sub]+"\"")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}
