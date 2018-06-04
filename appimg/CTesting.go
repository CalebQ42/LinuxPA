package appimg

// Thanks to probono for the majority of this C code: https://discourse.appimage.org/t/accessing-appimage-desktop-file/173/4?u=calebq42

/*
#cgo CFLAGS: -I/usr/lib
#cgo LDFLAGS: -L/usr/lib/libappimage.so -lappimage
#cgo pkg-config: glib-2.0

#include <appimage/appimage.h>
#include <stdlib.h>
#include <glib.h>
#include <strings.h>

bool extract_desktop(char* appimageLoc,char* extractLoc){
  char** files = appimage_list_files(appimageLoc);
  g_autofree gchar *desktop_file = NULL;
  gchar *extracted_desktop_file = extractLoc;
  int i = 0;
  for (; files[i] != NULL ; i++) {
      // g_debug("AppImage file: %s", files[i]);
      if (g_str_has_suffix (files[i],".desktop")) {
          desktop_file = files[i];
          g_debug("AppImage desktop file: %s", desktop_file);
          break;
      }
  }
  if(desktop_file == NULL) {
      g_debug("AppImage desktop file not found");
      appimage_string_list_free(files);
      return FALSE;
  }
  appimage_extract_file_following_symlinks(appimageLoc, desktop_file,extracted_desktop_file);
  appimage_string_list_free(files);
	return TRUE;
}

bool extract_file(char* appimageLoc,char* filename,char* extractLoc){
  char** files = appimage_list_files(appimageLoc);
  g_autofree gchar *found_file = NULL;
  gchar *extracted_file = extractLoc;
  int i = 0;
  for (; files[i] != NULL ; i++) {
      g_debug("AppImage file: %s", files[i]);
      if (strcmp(files[i],filename)==0) {
          found_file = files[i];
          g_debug("FileFound: %s", found_file);
      }
  }
  if(found_file == NULL) {
      g_debug("filenotfound");
      appimage_string_list_free(files);
      return FALSE;
  }
  appimage_extract_file_following_symlinks(appimageLoc, found_file,extracted_file);
  appimage_string_list_free(files);
	return TRUE;
}
*/
import "C"
import (
	"errors"
	"os"
	"unsafe"
)

func GetDesktopFile(loc string) (*os.File, error) {
	os.Remove("/tmp/my.desktop")
	var locTmp *C.char = C.CString(loc)
	defer C.free(unsafe.Pointer(locTmp))
	var extractLoc *C.char = C.CString("/tmp/my.desktop")
	defer C.free(unsafe.Pointer(extractLoc))
	var out C.bool = C.extract_desktop(locTmp, extractLoc)
	if out == false {
		return nil, errors.New("Desktop File Not Found!")
	}
	return os.Open("/tmp/my.desktop")
}
