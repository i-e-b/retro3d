package main

import (
	"fmt"
	"log"
	"retro3d/pkg"
	"syscall"
	"time"
	"unsafe"
)

var rend *pkg.Renderer
var lastFrame time.Time

func main() {
	rend = pkg.NewRenderer(800,600)
	rend.Update(0)
	rend.Update(1)
	go drawLoop()

	className := "testClass"

	instance, err := getModuleHandle()
	if err != nil {
		log.Println(err)
		return
	}

	cursor, err := loadCursorResource(cIDC_ARROW)
	if err != nil {
		log.Println(err)
		return
	}

	fn :=WindowsMessageHandler

	wcx := tWNDCLASSEXW{
		wndProc:    syscall.NewCallback(fn),
		instance:   instance,
		cursor:     cursor,
		background: cCOLOR_WINDOW + 1,
		className:  syscall.StringToUTF16Ptr(className),
	}
	wcx.size = uint32(unsafe.Sizeof(wcx))

	if _, err = registerClassEx(&wcx); err != nil {
		log.Println(err)
		return
	}

	_, err = createWindow(
		className,
		"Hello window",
		cWS_VISIBLE|cWS_OVERLAPPEDWINDOW,
		cSW_USE_DEFAULT,
		cSW_USE_DEFAULT,
		800,//cSW_USE_DEFAULT,
		600,//cSW_USE_DEFAULT,
		0,
		0,
		instance,
	)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		msg := tMSG{}
		gotMessage, err := getMessage(&msg, 0, 0, 0)
		if err != nil {
			log.Println(err)
			return
		}

		if gotMessage {
			translateMessage(&msg)
			dispatchMessage(&msg)
		} else {
			break
		}
	}
}

var lastHwnd syscall.Handle
var hwndIsSet bool = false
func drawLoop() {
	for {
		if !hwndIsSet{
			time.Sleep(250*time.Millisecond)
			continue
		}

		fTime := time.Since(lastFrame).Milliseconds()
		if fTime > 500 {fTime=500} // prevent jumps on hesitation

		if fTime > 15 { // don't go too fast
			go rend.Update(fTime) // update is parallel, frames are double-buffered
			DrawBitsIntoWindow(lastHwnd)
			lastFrame = time.Now()
		}

		time.Sleep(5*time.Millisecond)
	}
}


// WindowsMessageHandler is the ---------------
//             MAIN WIN32 EVENT LOOP
// --------------------------------------------
func WindowsMessageHandler (hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	lastHwnd = hwnd
	hwndIsSet = true
	switch msg {
	case cWM_CLOSE:
		destroyWindow(hwnd)
	case cWM_DESTROY:
		postQuitMessage(0)
	case cWM_ERASEBKGND:
		// do nothing, we will overdraw everything
		return 1
	case cWM_PAINT:
		paint:=tPAINTSTRUCT{}
		beginPaint(hwnd, &paint)
		endPaint(hwnd, &paint)
		return 1//ret
	default:
		ret := defWindowProc(hwnd, msg, wparam, lparam)
		return ret
	}
	return 0
}

func DrawBitsIntoWindow(hwnd syscall.Handle) {
	rc:= tRECT{}
	getClientRect(hwnd, &rc)

	width := rc.right
	height := rc.bottom
	if width < 1 || height < 1 {
		fmt.Println("error: invalid rectangle size")
		return
	}

	hdc := getDC(hwnd)
	defer releaseDC(hwnd, hdc)

	frame := rend.RenderFrame()

	// directly copy byte values to device
	// This in *kinda* working, but not 100%. Probably needs some invalidation or such-like in the event loop

	bmpWidth := int32(frame.Width)
	bmpHeight := int32(frame.Height)

	minHeight := min(bmpHeight, height)
	minWidth := min(bmpWidth, width)

	myBMInfo := tBITMAPINFO{}
	myBMInfo.bmiHeader = tBITMAPINFOHEADER{
		biWidth :    bmpWidth,
		biHeight :   bmpHeight,
		biPlanes :   1,
		biBitCount : 32,
	}
	myBMInfo.bmiHeader.biSize = int32(unsafe.Sizeof(myBMInfo))

	// SetDIBitsToDevice seems to be limited in how big a region it can copy.
	// Might need to do in 512x512 max chunks
	setDIBitsToDevice(hdc, 0,0, minWidth, minHeight, 0, 0, 0, uint32(minHeight), frame.GetBufferPointer(), &myBMInfo, DIB_RGB_COLORS)

}

func min(a,b int32) int32 {
	if a < b {return a}
	return b
}

