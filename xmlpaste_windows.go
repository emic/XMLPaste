package main

import (
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernelMod    = syscall.NewLazyDLL("kernel32")
	globalSize   = kernelMod.NewProc("GlobalSize")
	globalLock   = kernelMod.NewProc("GlobalLock")
	globalUnlock = kernelMod.NewProc("GlobalUnlock")

	userMod                = syscall.NewLazyDLL("user32")
	openClipboard          = userMod.NewProc("OpenClipboard")
	closeClipboard         = userMod.NewProc("CloseClipboard")
	getClipboardData       = userMod.NewProc("GetClipboardData")
	enumClipboardFormats   = userMod.NewProc("EnumClipboardFormats")
	getClipboardFormatName = userMod.NewProc("GetClipboardFormatNameA")
)

func getClipboard() (string, error) {
	var ret uintptr
	var err error

	ret, _, err = openClipboard.Call(0)
	defer closeClipboard.Call()
	if ret == 0 {
		return "", err
	}

	clipboardFormat, _, err := enumClipboardFormats.Call(0)
	if clipboardFormat != 13 { // 13: Custom Menu
		maxCount := 9
		buf := make([]uint16, maxCount)
		r, _, _ := getClipboardFormatName.Call(
			uintptr(clipboardFormat),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(maxCount))
		if r > 0 {
			clipboardTypes := map[string]string{
				"Mac-XMTB": "Table",
				"Mac-XMFD": "Field",
				"Mac-XMSC": "Script",
				"Mac-XMSS": "Script Step",
				"Mac-XMFN": "Custom Function",
				"Mac-XMLO": "Layout Object (.fp7)",
				"Mac-XML2": "Layout Object (.fmp12)",
				"Mac-XMVL": "Value List",
				"Mac-":     "Theme",
				"Mac-XMTH": "Theme (2024)",
			}
			clipboardType := string((*[1 << 4]byte)(unsafe.Pointer(&buf[0]))[0 : maxCount-1])
			if _, ex := clipboardTypes[clipboardType]; !ex && !(clipboardType[0:4] == "Mac-" && clipboardType[0:6] != "Mac-XM") {
				return "", err
			}
		} else {
			return "", err
		}
	}

	ret, _, err = getClipboardData.Call(clipboardFormat)
	if ret == 0 {
		return "", err
	}

	size, _, err := globalSize.Call(ret)
	if size == 0 {
		return "", err
	}

	l, _, err := globalLock.Call(ret)
	defer globalUnlock.Call(ret)
	if l == 0 {
		return "", err
	}

	if clipboardFormat == 13 {
		// for custome menus
		xml := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[0:])
		if strings.Index(xml, "<?xml version=\"1.0\" encoding=\"utf-16\"?><FMObjectTransfer ") == -1 {
			xml = ""
		} else if strings.HasSuffix(xml, "</FMObjectTransfer") {
			// for FileMaker Pro 17 Advanced for Windows
			xml = xml + ">"
		}
		return xml, nil
	} else {
		start := 4
		return string((*[1 << 20]byte)(unsafe.Pointer(l))[start:size]), nil
	}
}
