//go:build linux

package main

import (
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	row    uint16
	col    uint16
	xpixel uint16
	ypixel uint16
}

func terminalLines() int {
	ws := winsize{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
	if errno == 0 && ws.row > 0 {
		return int(ws.row)
	}
	return terminalLinesFallback()
}
