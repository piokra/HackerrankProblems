package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

func imin(xs ...int32) int32 {
	mx := int32(math.MaxInt32)
	for _, x := range xs {
		if mx > x {
			mx = x
		}
	}
	return mx
}

func loopShiftToXY(loop, shift, h, w int32) (int32, int32) {
	loopLength := 2*h + 2*w - 4*(2*loop+1)
	shift %= loopLength
	loopHeight := h - 2*loop
	loopWidth := w - 2*loop

	if shift < loopHeight {
		return loop, loop + shift
	}

	shift -= loopHeight - 1

	if shift < loopWidth {
		return loop + shift, loop + loopHeight - 1
	}

	shift -= loopWidth - 1

	if shift < loopHeight {
		return loop + loopWidth - 1, loop + loopHeight - 1 - shift
	}

	shift -= loopHeight - 1

	return loop + loopWidth - 1 - shift, loop

}

func xyToLoopShift(x, y, h, w int32) (int32, int32) {
	loop := imin(imin(x, y), imin(w-x-1, y), imin(x, h-y-1), imin(w-x-1, h-y-1))
	loopHeight := h - 2*loop
	loopWidth := w - 2*loop

	x -= loop
	y -= loop

	if x == 0 {
		return loop, y
	}

	if y == loopHeight - 1 {
		return loop, loopHeight - 1 + x
	}

	if x == loopWidth - 1 {
		return loop, loopHeight + loopWidth - 2 + (loopHeight-1-y)
	}

	loopLength := 2*loopWidth + 2*loopHeight - 4
	return loop, loopLength - x
}

func getShiftedValue(matrix [][]int32, x, y, r int32) int32 {
	h, w := int32(len(matrix)), int32(len(matrix[0]))
	loop := imin(imin(x, y), imin(w-x-1, y), imin(x, h-y-1), imin(w-x-1, h-y-1))
	loop, shift := xyToLoopShift(x, y, h, w)
	loopLength := 2*h + 2*w - 4*(2*loop+1)

	shift += loopLength - (r % loopLength)
	nx, ny := loopShiftToXY(loop, shift, h, w)
	return matrix[ny][nx]
}

// Complete the matrixRotation function below.
func matrixRotation(matrix [][]int32, r int32) {
	w := len(matrix[0])
	for y, row := range matrix {
		for x := range row {
			v := getShiftedValue(matrix, int32(x), int32(y), r)
			if x + 1 != w {
				fmt.Print(v, " ")
			} else {
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	mnr := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	mTemp, err := strconv.ParseInt(mnr[0], 10, 64)
	checkError(err)
	m := int32(mTemp)

	nTemp, err := strconv.ParseInt(mnr[1], 10, 64)
	checkError(err)
	n := int32(nTemp)

	rTemp, err := strconv.ParseInt(mnr[2], 10, 64)
	checkError(err)
	r := int32(rTemp)

	var matrix [][]int32
	for i := 0; i < int(m); i++ {
		matrixRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var matrixRow []int32
		for _, matrixRowItem := range matrixRowTemp {
			matrixItemTemp, err := strconv.ParseInt(matrixRowItem, 10, 64)
			checkError(err)
			matrixItem := int32(matrixItemTemp)
			matrixRow = append(matrixRow, matrixItem)
		}

		if len(matrixRow) != int(n) {
			panic("Bad input")
		}

		matrix = append(matrix, matrixRow)
	}

	matrixRotation(matrix, r)
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
