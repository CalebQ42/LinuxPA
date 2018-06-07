package appimg

// Thanks to probono for the majority of this C code: https://discourse.appimage.org/t/accessing-appimage-desktop-file/173/4?u=calebq42

/*
#cgo CFLAGS: -I/usr/lib
#cgo LDFLAGS: -L/usr/lib/libappimage.so -lappimage

#include <appimage/appimage.h>
#include <stdlib.h>

int char_length(char** in){
	int i = 0;
	for (; in[i] != NULL ; i++);
	return i;
}
*/
import "C"
import (
	"errors"
	"os"
	"strings"
	"unsafe"
)

var (
	//ErrNotAppImage return if the given location does not point to an AppImage
	ErrNotAppImage = errors.New("file is not an appimage")
)

//Thanks to https://stackoverflow.com/questions/36188649/cgo-char-to-slice-string for char** to string slice code :)

//GetDesktopFile tries to find the desktop file inside the appimage at appimageLoc and extract it to extractLoc.
//extractLoc is deleted first to provent conflicts.
func GetDesktopFile(appimageLoc string, extractLoc string) (*os.File, error) {
	if !strings.HasSuffix(appimageLoc, ".AppImage") {
		return nil, ErrNotAppImage
	}
	err := os.Remove(extractLoc)
	if err != os.ErrNotExist {
		return nil, err
	}
	cextractLoc := C.CString(extractLoc)
	defer C.free(unsafe.Pointer(cextractLoc))
	cloc := C.CString(appimageLoc)
	defer C.free(unsafe.Pointer(cloc))
	cfiles := C.appimage_list_files(cloc)
	defer C.appimage_string_list_free(cfiles)
	cfilesLength := C.char_length(cfiles)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(cfiles))[:cfilesLength:cfilesLength]
	for _, v := range tmpslice {
		tmp := C.GoString(v)
		if strings.HasSuffix(tmp, ".desktop") {
			C.appimage_extract_file_following_symlinks(cloc, v, cextractLoc)
			break
		}
	}
	return os.Open(extractLoc)
}

//ExtractFile tries to extract the file at fileLoc to extractLoc from the appimage at appimageLoc.
//extractLoc is deleted first to provent conflicts.
func ExtractFile(appimageLoc string, fileLoc string, extractLoc string) (*os.File, error) {
	if !strings.HasSuffix(appimageLoc, ".AppImage") {
		return nil, ErrNotAppImage
	}
	err := os.Remove(extractLoc)
	if err != os.ErrNotExist {
		return nil, err
	}
	cextractLoc := C.CString(extractLoc)
	defer C.free(unsafe.Pointer(cextractLoc))
	cloc := C.CString(appimageLoc)
	defer C.free(unsafe.Pointer(cloc))
	cfileLoc := C.CString(fileLoc)
	defer C.free(unsafe.Pointer(cfileLoc))
	C.appimage_extract_file_following_symlinks(cloc, cfileLoc, cextractLoc)
	return os.Open(extractLoc)
}
