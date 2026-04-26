package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/skip2/go-qrcode"
)

const (
	qrANSIReset     = "\x1b[0m"
	qrANSIBlackBG   = "\x1b[40m"
	qrANSIWhiteBG   = "\x1b[47m"
	qrANSIBlackOnW  = "\x1b[30;47m"
	qrLargeCell     = "  "
	qrCompactWhite  = ' '
	qrCompactTop    = '▀'
	qrCompactBottom = '▄'
	qrCompactFull   = '█'
)

func renderTerminalQRCode(rawURL string) {
	code, err := qrcode.New(rawURL, qrcode.Medium)
	if err != nil {
		fmt.Printf("终端二维码渲染失败，使用下方地址登录: %v\n", err)
		return
	}

	bitmap := code.Bitmap()
	if useCompactQRCode(bitmap) {
		renderCompactQRCode(bitmap)
		return
	}
	renderLargeQRCode(bitmap)
}

func useCompactQRCode(bitmap [][]bool) bool {
	if len(bitmap) == 0 || len(bitmap[0]) == 0 {
		return false
	}

	cols, rows, ok := readTerminalSize()
	if !ok {
		return true
	}

	largeWidth := len(bitmap[0]) * len(qrLargeCell)
	largeHeight := len(bitmap)
	return cols < largeWidth+2 || rows < largeHeight+2
}

func renderLargeQRCode(bitmap [][]bool) {
	for _, row := range bitmap {
		fmt.Println(renderQRCodeRow(row))
	}
	fmt.Print(qrANSIReset)
}

func renderQRCodeRow(row []bool) string {
	if len(row) == 0 {
		return ""
	}

	var line strings.Builder
	current := row[0]
	line.WriteString(qrANSIColor(current))
	line.WriteString(qrLargeCell)

	for i := 1; i < len(row); i++ {
		cell := row[i]
		if cell != current {
			line.WriteString(qrANSIColor(cell))
			current = cell
		}
		line.WriteString(qrLargeCell)
	}

	line.WriteString(qrANSIReset)
	return line.String()
}

func renderCompactQRCode(bitmap [][]bool) {
	for y := 0; y < len(bitmap); y += 2 {
		fmt.Println(renderCompactRow(bitmap, y))
	}
	fmt.Print(qrANSIReset)
}

func renderCompactRow(bitmap [][]bool, y int) string {
	top := bitmap[y]
	var bottom []bool
	if y+1 < len(bitmap) {
		bottom = bitmap[y+1]
	}

	var line strings.Builder
	line.WriteString(qrANSIBlackOnW)
	for x := 0; x < len(top); x++ {
		topCell := top[x]
		bottomCell := false
		if len(bottom) > x {
			bottomCell = bottom[x]
		}
		line.WriteRune(compactQRCodeRune(topCell, bottomCell))
	}
	line.WriteString(qrANSIReset)
	return line.String()
}

func compactQRCodeRune(top, bottom bool) rune {
	switch {
	case top && bottom:
		return qrCompactFull
	case top:
		return qrCompactTop
	case bottom:
		return qrCompactBottom
	default:
		return qrCompactWhite
	}
}

func qrANSIColor(cell bool) string {
	if cell {
		return qrANSIBlackBG
	}
	return qrANSIWhiteBG
}

func readTerminalSize() (cols, rows int, ok bool) {
	if cols, rows, ok = readTerminalSizeFromEnv(); ok {
		return cols, rows, true
	}

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}
	var parsedRows, parsedCols int
	if _, err := fmt.Sscanf(string(out), "%d %d", &parsedRows, &parsedCols); err != nil {
		return 0, 0, false
	}
	if parsedCols <= 0 || parsedRows <= 0 {
		return 0, 0, false
	}
	return parsedCols, parsedRows, true
}

func readTerminalSizeFromEnv() (cols, rows int, ok bool) {
	var parsedCols, parsedRows int
	if _, err := fmt.Sscanf(os.Getenv("COLUMNS"), "%d", &parsedCols); err != nil || parsedCols <= 0 {
		return 0, 0, false
	}
	if _, err := fmt.Sscanf(os.Getenv("LINES"), "%d", &parsedRows); err != nil || parsedRows <= 0 {
		return 0, 0, false
	}
	return parsedCols, parsedRows, true
}
