package apps

import (
	"bytes"
	"debug/elf"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/probonopd/go-appimage/src/goappimage"
)

type Exe struct {
	path     string
	name     string
	args     string //Obtained from .desktop files
	appimage bool
	script   bool
}

func ProcessExe(file string) (Exe, error) {
	var err error
	var e Exe
	if !path.IsAbs(file) {
		var wd string
		wd, err = os.Getwd()
		if err != nil {
			return e, err
		}
		file = path.Join(wd, file)
	}
	file = path.Clean(file)
	e.path = file
	e.name = path.Base(file)
	if strings.Contains(e.name, ".") {
		e.name = e.name[:strings.LastIndex(e.name, ".")]
	}
	//Check if file exists/is openable
	fil, err := os.Open(file)
	if err != nil {
		return e, err
	}
	//If it's a windows executable, we don't need to process further
	if e.IsWine() {
		return e, nil
	}
	//Then Check for a shebang (script file)
	check := make([]byte, 2)
	_, err = fil.Read(check)
	if err != nil {
		return e, e.CheckSetExec()
	}
	if bytes.Equal(check, []byte("#!")) {
		e.script = true
		return e, e.CheckSetExec()
	}
	//Then we check if the file is an ELF file that's executable
	elfFile, err := elf.Open(file)
	if err != nil {
		return e, e.CheckSetExec()
	}
	if elfFile.Type != elf.ET_EXEC {
		return e, errors.New("file is not executable elf")
	}
	//Lastly, check if it's an AppImage so we can potentially enable advanced features
	_, err = goappimage.NewAppImage(file)
	e.appimage = err == nil
	return e, nil
}

func (e Exe) String() string {
	out := "Executable: " + e.name + " at path: " + e.path
	if e.args != "" {
		out += " has args: \"" + e.args + "\""
	}
	if e.appimage {
		out += " is appimage"
	}
	if e.script {
		out += " is script"
	}
	return out
}

func (e Exe) IsWine() bool {
	return strings.HasSuffix(e.path, ".exe")
}

func (e Exe) Cmd(commonSh, wine string, fromRoot bool) (cmd *exec.Cmd) {
	if commonSh != "" {
		cmd = exec.Command(commonSh)
		if e.IsWine() {
			cmd.Args = append(cmd.Args, wine)
		}
		cmd.Args = append(cmd.Args, e.path)
		env := os.Environ()
		wd, _ := os.Getwd()
		env = append(env, "ROOT="+wd, "APPNAME="+e.name, "FILENAME="+e.path)
		cmd.Env = env
	} else if e.IsWine() {
		cmd = exec.Command(wine, e.path)
	} else {
		cmd = exec.Command(e.path)
	}
	if !fromRoot {
		cmd.Dir = path.Dir(e.path)
	}
	cmd.Args = append(cmd.Args, strings.Split(e.args, " ")...)
	return
}

func (e Exe) CheckSetExec() error {
	if e.HasExecPerm() {
		return nil
	}
	return e.SetExecPerm()
}

func (e Exe) HasExecPerm() bool {
	stat, err := os.Stat(e.path)
	if err != nil {
		return false
	}
	perm := stat.Mode()
	sysStat := stat.Sys().(*syscall.Stat_t)
	if os.Getuid() == int(sysStat.Uid) {
		if perm&0100 == 0100 {
			return true
		}
	}
	groups, err := os.Getgroups()
	if err != nil {
		groups = []int{os.Getgid()}
	}
	for i := range groups {
		if groups[i] == int(sysStat.Gid) {
			return perm&0010 == 0010
		}
	}
	return perm&0001 == 0001
}

func (e Exe) SetExecPerm() error {
	stat, err := os.Stat(e.path)
	if err != nil {
		return err
	}
	perm := stat.Mode()
	sysStat := stat.Sys().(*syscall.Stat_t)
	if os.Getuid() == int(sysStat.Uid) {
		return os.Chmod(e.path, perm&0100)
	}
	groups, err := os.Getgroups()
	if err != nil {
		groups = []int{os.Getgid()}
	}
	for i := range groups {
		if groups[i] == int(sysStat.Gid) {
			return os.Chmod(e.path, perm&0010)
		}
	}
	return errors.New("not user or group owner and for security not setting other exec permission")
}