//go:build !linux

package main

func terminalLines() int {
	return terminalLinesFallback()
}
